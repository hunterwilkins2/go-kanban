package config

type Config struct {
	Port      int
	DSN       string
	Version   string
	Env       string
	HotReload bool
}
