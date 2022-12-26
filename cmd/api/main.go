package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/kl09/counter-api/api"
	"github.com/kl09/counter-api/internal/counter"
	"github.com/kl09/counter-api/internal/events"
	"github.com/oklog/run"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	logger.Log("msg", "starting api service")

	fs := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	{
		fs.String("http-addr", ":8080", "Address to listen for System API")
	}

	if err := viper.BindPFlags(fs); err != nil {
		logger.Log("msg", "failed bind pflags", "error", err)
		os.Exit(1)
	}

	evCh := events.NewEventCh(logger)

	router := api.NewRouter(
		api.NewHandlerConfig(
			api.WithLogger(log.WithPrefix(logger, "component", "api.Router")),
		),
		counter.NewStats(evCh),
		evCh,
	)

	apiServer := http.Server{
		Addr:        viper.GetString("http-addr"),
		ReadTimeout: 20 * time.Second,
		IdleTimeout: 30 * time.Second,
		Handler:     router.Handler(),
	}

	// This ctx is used to terminate the go routine that listens for system signals, nothing more.
	ctx, cancel := context.WithCancel(context.Background())

	var g run.Group
	{
		g.Add(func() error {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			select {
			case <-ctx.Done():
				return nil
			case <-sig:
				level.Info(logger).Log("msg", "terminating...")
				return nil
			}
		}, func(_ error) {
			level.Info(logger).Log("msg", "program was interrupted")
			cancel()
		})
	}
	{
		g.Add(func() error {
			level.Info(logger).Log("started server for addr: ", apiServer.Addr)
			return apiServer.ListenAndServe()
		}, func(_ error) {
			level.Info(logger).Log("msg", "program was interrupted")
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			shErr := apiServer.Shutdown(ctx)

			if shErr != nil {
				level.Info(logger).Log("msg", "API server shut down with error", "err", shErr)
			}
		})
	}
	{
		g.Add(func() error {
			level.Info(logger).Log("msg", "started event listener")
			return evCh.Listen(ctx)
		}, func(_ error) {
			level.Info(logger).Log("msg", "program was interrupted")
			cancel()
		})
	}

	err := g.Run()
	if err != nil {
		level.Error(logger).Log("msg", "actors stopped gracefully", "err", err)
		os.Exit(1)
	}

	os.Exit(0)
}
