package hook

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/murtaza-u/amify/internal/conf"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const port = ":5748"

// Hook represents a webhook object.
type Hook struct {
	conf      conf.C
	startedAt time.Time
}

// New initializes a webhook object with the provided configuration.
func New(c conf.C) (*Hook, error) {
	return &Hook{
		conf:      c,
		startedAt: time.Now(),
	}, nil
}

// Listen starts the webhook API server.
func (h Hook) Listen() {
	e := echo.New()

	// configure logger
	slog.SetDefault(h.logger())

	// setup basic auth middleware, if enabled
	var middlewares []echo.MiddlewareFunc
	if h.conf.Hook.Auth.Enable {
		middlewares = append(middlewares, middleware.BasicAuth(h.basicAuth))
	}

	e.POST("/hook", h.serve, middlewares...)
	e.GET("/health", h.health)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := e.Start(port)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Log(context.Background(), slog.LevelInfo, "shutting down")
				return
			}
			slog.LogAttrs(
				context.Background(),
				slog.LevelError,
				"server terminated",
				slog.String("error", err.Error()),
			)
			stop()
		}
	}()

	// interrupt received
	<-ctx.Done()

	// graceful termination
	ctx, cancel := context.WithCancel(context.Background())
	if h.conf.Hook.TerminationGracePeriod != 0 {
		ctx, cancel = context.WithTimeout(ctx, h.conf.Hook.TerminationGracePeriod)
	}
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"forcefully shutting down",
			slog.String("error", err.Error()),
		)
	}

	wg.Wait()
}
