package conf

// C contains all configuration data that can be passed to the webhook.
type C struct {
	// Hook contains the webhook configuration.
	Hook Hook `koanf:"hook"`
	// Ntfy contains the ntfy server configuration.
	Ntfy Ntfy `koanf:"ntfy"`
}

// Hook contains all configuration related to the webhook.
type Hook struct {
	// ListenAddr is the address the webhook should listen on.
	//
	// Default: ":8080"
	ListenAddr string `koanf:"listenAddr"`
	// Auth contains the configuration for securing the webhook endpoint with
	// HTTP basic authentication.
	Auth Auth `koanf:"auth"`
	// Log contains the configuration for the log format and level.
	Log Log `koanf:"log"`
}

// Ntfy contains all configuration related to ntfy.
type Ntfy struct {
	// BaseURL is the ntfy server's base URL. For example: https://ntfy.sh
	//
	// Required.
	BaseURL string `koanf:"baseUrl"`
	// Auth contains the configuration for authenticating with the ntfy server.
	Auth Auth `koanf:"auth"`
	// Notification contains the configuration for notification messages.
	Notification Notification `koanf:"notification"`
}

// Auth contains HTTP basic authentication configuration.
type Auth struct {
	// Enable HTTP basic authentication.
	//
	// For the webhook, enabling this will protect the hook endpoint with HTTP
	// basic authentication.
	//
	// For ntfy, enabling this will allow the webhook to use HTTP basic auth
	// credentials to authenticate with the ntfy server.
	//
	// Default: false
	Enable bool `koanf:"enable"`
	// Username is the HTTP basic auth username. Required if `Enable` is true.
	Username string `koanf:"username"`
	// Password is the HTTP basic auth password. Required if `Enable` is true.
	Password string `koanf:"password"`
}

// Log contains configuration for the webhook logger.
type Log struct {
	// Level is the log level.
	// Possible values: "debug", "info", "warn", "error".
	//
	// Default: "info"
	Level string `koanf:"level"`
	// Format specifies the format of the logs.
	// Possible values: "text", "json"
	//
	// Default: "text"
	Format string `koanf:"format"`
}

// Notification contains configuration for notification messages.
type Notification struct {
	// Topic can be a hardcoded string or a gval expression that evaluates to a
	// string. For example: "alertmanager"
	//
	// Required.
	Topic StringExpr `koanf:"topic"`
	// Priority can be a hardcoded string or a gval expression that evaluates
	// to a string.
	// For example: Status == "firing" ? "urgent" : "default"
	//
	// Reference: https://docs.ntfy.sh/publish/#message-priority
	//
	// Default: "default"
	Priority StringExpr `koanf:"priority"`
	// Tags to be included in the notification. Optional.
	// Emoji shortcode reference: https://docs.ntfy.sh/emojis/
	Tags []Tag `koanf:"tags"`
	// Title of the notification. Required.
	Title Template `koanf:"title"`
	// Description of the notification. Required.
	Description Template `koanf:"description"`
}

// Tag represents a tag to be included with the notification.
type Tag struct {
	// Tag can be an emoji shortcode or any contextual data.
	Tag string `koanf:"tag"`
	// Condition is a gval expression. The tag is included in the notification
	// only if the condition evaluates to true or is empty.
	Condition Expr `koanf:"condition"`
}
