package hook

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Hook) serve(c echo.Context) error {
	c.NoContent(http.StatusAccepted)
	return nil
}
