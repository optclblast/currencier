package config

type Config struct {
	Common CommonConfig
	Rest   RestConfig
	Cache  Cache
}

type CommonConfig struct {
	Level string
}

type RestConfig struct {
	Port string
	TLS  bool
}

type Cache struct {
	URL    string
	User   string
	Secret string
}
