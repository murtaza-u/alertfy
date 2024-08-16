package ntfy

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/murtaza-u/amify/internal/conf"
)

// RequestData contains the data used to create an HTTP request.
type RequestData struct {
	Notification Data
	BasicAuth    conf.Auth
}

// NewRequest creates a new HTTP request to the ntfy server, including all the
// details about the notification message.
func NewRequest(ctx context.Context, data RequestData) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		data.Notification.URL,
		strings.NewReader(data.Notification.Description),
	)
	if err != nil {
		return nil, fmt.Errorf("new http request: %w", err)
	}

	if data.BasicAuth.Enable {
		uname := data.BasicAuth.Username
		pswd := data.BasicAuth.Password
		req.SetBasicAuth(uname, pswd)
	}

	if data.Notification.Title != "" {
		req.Header.Set("X-Title", data.Notification.Title)
	}
	if data.Notification.Tags != "" {
		req.Header.Set("X-Tags", data.Notification.Tags)
	}
	req.Header.Set("X-Priority", data.Notification.Priority)

	return req, nil
}
