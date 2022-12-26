package api

import (
	"github.com/go-kit/log"
)

// HandlerConfig is a config for NewHTTPHandler function.
type HandlerConfig struct {
	// logger is logger used by Handlers.
	logger log.Logger
}

// NewHandlerConfig creates a new config for NewHTTPHandler function.
func NewHandlerConfig(options ...HandlerOption) HandlerConfig {
	c := HandlerConfig{
		logger: log.NewNopLogger(),
	}

	for _, opt := range options {
		opt(&c)
	}

	return c
}

// HandlerOption configures HandlerConfig.
type HandlerOption func(config *HandlerConfig)

// WithLogger configures a logger to debug Handlers.
func WithLogger(l log.Logger) func(config *HandlerConfig) {
	return func(c *HandlerConfig) {
		c.logger = l
	}
}
