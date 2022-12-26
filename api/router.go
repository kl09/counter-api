package api

import (
	"net/http"

	"github.com/gorilla/mux"
	domain "github.com/kl09/counter-api"
)

type Router struct {
	config           HandlerConfig
	statsCounter     domain.StatsCounter
	eventWsConnector domain.EventWSConnector
}

func NewRouter(config HandlerConfig, statsCounter domain.StatsCounter, eventWsConnector domain.EventWSConnector) *Router {
	return &Router{
		config:           config,
		statsCounter:     statsCounter,
		eventWsConnector: eventWsConnector,
	}
}

func (ro *Router) Handler() http.Handler {
	r := mux.NewRouter()

	r.Methods("GET").Path("/v1/subscribe").HandlerFunc(ro.subscriberHandler)

	r.Methods("GET").Path("/v1/{key}").HandlerFunc(ro.getKeyHandler)
	r.Methods("POST").Path("/v1/{key}/increment").HandlerFunc(ro.incrementHandler)
	r.Methods("POST").Path("/v1/{key}/decrement").HandlerFunc(ro.decrementHandler)
	r.Methods("POST").Path("/v1/{key}/reset").HandlerFunc(ro.resetHandler)

	if ro.config.logger != nil {
		r.Use(newLogger(ro.config.logger))
	}

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}
