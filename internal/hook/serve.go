package hook

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type request struct {
	Receiver string  `json:"receiver"`
	Status   string  `json:"status"`
	Alerts   []alert `json:"alerts"`
}

type alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

func (h Hook) serve(c echo.Context) error {
	req := new(request)
	if err := c.Bind(req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.NoContent(http.StatusAccepted)
}
