package hook

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h Hook) health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
		"uptime": time.Since(h.startedAt).String(),
	})
}
