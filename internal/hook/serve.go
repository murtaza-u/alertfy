package hook

import (
	"log/slog"
	"net/http"

	"github.com/murtaza-u/alertfy/internal/alert"
	"github.com/murtaza-u/alertfy/internal/ntfy"

	"github.com/labstack/echo/v4"
)

type request struct {
	Receiver    string        `json:"receiver"`
	Status      string        `json:"status"`
	Alerts      []alert.Alert `json:"alerts"`
	ExternalURL string        `json:"externalURL"`
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

		err := h.forwardAlert(c, alert)
		if err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusAccepted)
}

func (h Hook) forwardAlert(c echo.Context, alert alert.Alert) error {
	p := ntfy.NewParser(h.conf.Ntfy)
	data := p.Parse(c.Request().Context(), alert)
	if data == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	req, err := ntfy.NewRequest(c.Request().Context(), ntfy.RequestData{
		Notification: *data,
		BasicAuth:    h.conf.Ntfy.Auth,
	})
	if err != nil {
		slog.LogAttrs(
			c.Request().Context(),
			slog.LevelError,
			"failed to create http request. Aborting",
			slog.String("fingerprint", alert.Fingerprint),
			slog.String("error", err.Error()),
		)
		return c.NoContent(http.StatusInternalServerError)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.LogAttrs(
			c.Request().Context(),
			slog.LevelError,
			"failed to forward request to ntfy server. Aborting",
			slog.String("fingerprint", alert.Fingerprint),
			slog.String("error", err.Error()),
		)
		return c.NoContent(http.StatusInternalServerError)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.LogAttrs(
			c.Request().Context(),
			slog.LevelError,
			"non-2XX status code received from ntfy server. Aborting",
			slog.String("fingerprint", alert.Fingerprint),
			slog.String("status", resp.Status),
		)
		return c.NoContent(http.StatusInternalServerError)
	}

	return nil
}
