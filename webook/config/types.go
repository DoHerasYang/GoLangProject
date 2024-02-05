package config

type config struct {
	DB      DBConfig
	Redis   RedisConfig
	GinHost GinConfig
}

type DBConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}

type GinConfig struct {
	Addr string
}
