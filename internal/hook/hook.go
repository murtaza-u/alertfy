package hook

import (
	"github.com/murtaza-u/amify/internal/conf"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Hook represents a webhook object.
type Hook struct {
	conf conf.C
}

// New initializes a webhook object with the provided configuration.
func New(c conf.C) (*Hook, error) {
	return &Hook{
		conf: c,
	}, nil
}

// Listen starts the webhook API server.
func (h Hook) Listen() error {
	e := echo.New()

	// setup basic auth middleware, if enabled
	if h.conf.Hook.Auth.Enable {
		e.Use(middleware.BasicAuth(h.basicAuth))
	}

	e.POST("/hook", h.serve)
	return e.Start(h.conf.Hook.ListenAddr)
}
