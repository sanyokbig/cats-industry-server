package config

import "github.com/caarlos0/env"

var EveConfig struct {
	ClientId    string `env:"CLIENT_ID"`
	SecretKey   string `env:"SECRET_KEY"`
	RedirectUri string `env:"REDIRECT_URI"`
}

var PostgresConfig struct {
	Host string `env:"POSTGRES_HOST"`
	Port string `env:"POSTGRES_PORT"`
	Db   string `env:"POSTGRES_DB"`
	User string `env:"POSTGRES_USER"`
	Pass string `env:"POSTGRES_PASS"`
}

var RedisConfig struct {
	Uri     string `env:"REDIS_URI"`
	DB      int    `env:"REDIS_DB"`
	Pass    string `env:"REDIS_PASS"`
	TTLDays int    `env:"REDIS_TTL_DAYS"`
}

func Parse() {
	err := env.Parse(&EveConfig)
	if err != nil {
		panic(err)
	}
	err = env.Parse(&PostgresConfig)
	if err != nil {
		panic(err)
	}
	err = env.Parse(&RedisConfig)
	if err != nil {
		panic(err)
	}
}
