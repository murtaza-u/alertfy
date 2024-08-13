package conf

// C contains all the configuration data that can be passed to the webhook.
type C struct {
	Hook hook `koanf:"hook"`
}

type hook struct {
	ListenAddr string `koan:"listenAddr"`
	Auth       auth   `koan:"auth"`
}

type auth struct {
	Username string `koan:"username"`
	Password string `koan:"password"`
}
