package hook

import (
	"net/http"

	"github.com/murtaza-u/amify/internal/conf"

	"github.com/labstack/echo/v4"
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
	e.POST("/hook", h.serve)
	return e.Start(h.conf.Hook.ListenAddr)
}

func (h Hook) serve(c echo.Context) error {
	c.NoContent(http.StatusAccepted)
	return nil
}
