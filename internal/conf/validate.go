package conf

import (
	"fmt"
	"net"
	"net/url"
)

// Validate validates the provided configuration.
func (c C) Validate() error {
	// hook
	if err := validateListenAddr(c.Hook.ListenAddr); err != nil {
		return fmt.Errorf("invalid `hook.listenAddr`: %q", c.Hook.ListenAddr)
	}
	if err := validateAuth(c.Hook.Auth); err != nil {
		return fmt.Errorf("`hook.auth`: %w", err)
	}
	if err := validateLogLevel(c.Hook.Log.Level); err != nil {
		return fmt.Errorf("`c.hook.log.level`: %w", err)
	}
	if err := validateLogFormat(c.Hook.Log.Format); err != nil {
		return fmt.Errorf("`c.hook.log.format`: %w", err)
	}
	if c.Hook.TerminationGracePeriod < 0 {
		return fmt.Errorf("`hook.terminationGracePeriod` cannot be -ve")
	}

	// ntfy
	if _, err := url.Parse(c.Ntfy.BaseURL); err != nil {
		return fmt.Errorf("invalid `ntfy.baseUrl` %q: %w", c.Ntfy.BaseURL, err)
	}
	if err := validateAuth(c.Ntfy.Auth); err != nil {
		return fmt.Errorf("`ntfy.auth`: %w", err)
	}
	if c.Ntfy.Notification.Topic.Text == "" {
		return fmt.Errorf("`ntfy.notification.topic` cannot be empty")
	}
	if c.Ntfy.Notification.Title == nil {
		return fmt.Errorf("`ntfy.notification.title` cannot be empty")
	}
	if c.Ntfy.Notification.Description == nil {
		return fmt.Errorf("`ntfy.notification.description` cannot be empty")
	}

	return nil
}

func validateListenAddr(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	// validate the host part
	if host != "" {
		if ip := net.ParseIP(host); ip == nil {
			return fmt.Errorf("invalid ip address: %q", host)
		}
	}

	// validate the port part
	if _, err := net.LookupPort("tcp", port); err != nil {
		return fmt.Errorf("invalid port %q: %w", port, err)
	}

	return nil
}

func validateAuth(auth Auth) error {
	if !auth.Enable {
		return nil
	}
	uname := auth.Username
	pswd := auth.Password
	if uname == "" {
		return fmt.Errorf("auth is enabled but `auth.username` is not set")
	}
	if pswd == "" {
		return fmt.Errorf("auth is enabled but `auth.password` is not set")
	}
	return nil
}

func validateLogLevel(level string) error {
	switch level {
	case "debug":
	case "info":
	case "warn":
	case "error":
	default:
		return fmt.Errorf("invalid value for `hook.log.level`: %q", level)
	}
	return nil
}

func validateLogFormat(format string) error {
	switch format {
	case "text":
	case "json":
	default:
		return fmt.Errorf("invalid value for `hook.log.format`: %q", format)
	}
	return nil
}
