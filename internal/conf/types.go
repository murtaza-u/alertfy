package conf

// C contains all the configuration data that can be passed to the webhook.
type C struct {
	Hook hook `koanf:"hook"`
	Ntfy ntfy `kaonf:"ntfy"`
}

type hook struct {
	ListenAddr string `koanf:"listenAddr"`
	Auth       auth   `koanf:"auth"`
	Log        log    `koanf:"log"`
}

type ntfy struct {
	BaseURL      string       `koanf:"baseUrl"`
	Auth         auth         `koanf:"auth"`
	Notification notification `koanf:"notification"`
}

type auth struct {
	Enable   bool   `koanf:"enable"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

type log struct {
	Level  string `koanf:"level"`
	Format string `koanf:"format"`
}

type notification struct {
	Topic    StringExpr `koanf:"topic"`
	Priority StringExpr `koanf:"priority"`
	Tags     []tag      `koanf:"tags"`
}

type tag struct {
	Tag       string `koanf:"tag"`
	Condition Expr   `koanf:"condition"`
}
