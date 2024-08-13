package hook

import (
	"crypto/subtle"

	"github.com/labstack/echo/v4"
)

func (h Hook) basicAuth(uname, pswd string, c echo.Context) (bool, error) {
	auth := h.conf.Hook.Auth
	uMatch := subtle.ConstantTimeCompare([]byte(uname), []byte(auth.Username))
	pMatch := subtle.ConstantTimeCompare([]byte(pswd), []byte(auth.Password))
	if uMatch == 1 && pMatch == 1 {
		return true, nil
	}
	return false, nil
}
