package conf

type Config struct {
	Hook Hook `koanf:"hook"`
}

type Hook struct {
	ListenAddr string `koan:"listenAddr"`
	Auth       Auth   `koan:"auth"`
}

type Auth struct {
	Username string `koan:"username"`
	Password string `koan:"password"`
}
