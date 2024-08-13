package hook

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type request struct {
	Receiver    string  `json:"receiver"`
	Status      string  `json:"status"`
	Alerts      []alert `json:"alerts"`
	ExternalURL string  `json:"externalURL"`
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
		slog.LogAttrs(
			c.Request().Context(),
			slog.LevelError,
			"failed to parse request body",
			slog.Int("status", http.StatusBadRequest),
		)
		return c.NoContent(http.StatusBadRequest)
	}

	if len(req.Alerts) == 0 {
		slog.LogAttrs(
			c.Request().Context(),
			slog.LevelWarn,
			"received request with zero alerts",
			slog.String("receiver", req.Receiver),
			slog.String("globalStatus", req.Status),
			slog.String("externalURL", req.ExternalURL),
		)
		return c.NoContent(http.StatusBadRequest)
	}

	for _, alert := range req.Alerts {
		slog.LogAttrs(
			c.Request().Context(),
			slog.LevelDebug,
			"received alert",
			slog.String("fingerprint", alert.Fingerprint),
			slog.String("receiver", req.Receiver),
			slog.String("status", alert.Status),
			slog.Time("startsAt", alert.StartsAt),
			slog.Time("endsAt", alert.EndsAt),
			slog.String("generatorURL", alert.GeneratorURL),
			slog.String("labels", formatLabels(alert.Labels)),
			slog.String("annotations", formatLabels(alert.Annotations)),
		)
	}

	return c.NoContent(http.StatusAccepted)
}
