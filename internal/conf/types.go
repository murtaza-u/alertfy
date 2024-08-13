package conf

// C contains all the configuration data that can be passed to the webhook.
type C struct {
	Hook hook `koanf:"hook"`
}

type hook struct {
	ListenAddr string `koan:"listenAddr"`
	Auth       auth   `koan:"auth"`
	Log        log    `koan:"log"`
}

type auth struct {
	Enable   bool   `koan:"enable"`
	Username string `koan:"username"`
	Password string `koan:"password"`
}

type log struct {
	Level  string `koan:"level"`
	Format string `koan:"format"`
}
