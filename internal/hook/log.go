package hook

import (
	"log/slog"
	"os"
)

func (h Hook) logger() *slog.Logger {
	var level slog.Level
	switch h.conf.Hook.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	opt := &slog.HandlerOptions{Level: level}
	if h.conf.Hook.Log.Format == "json" {
		return slog.New(slog.NewJSONHandler(os.Stderr, opt))
	}
	return slog.New(slog.NewTextHandler(os.Stderr, opt))
}
