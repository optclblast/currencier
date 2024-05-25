package config

type Config struct {
	Common CommonConfig
	Rest   RestConfig
	PG     PG
}

type CommonConfig struct {
	Level string
}

type RestConfig struct {
	Port string
	TLS  bool
}

type PG struct {
	URL     string
	PoolMax int
}
