package config

import "github.com/caarlos0/env"

var EveConfig struct {
	ClientId    string `env:"CLIENT_ID"`
	SecretKey   string `env:"SECRET_KEY"`
	RedirectUri string `env:"REDIRECT_URI"`
}

var MongoConfig struct {
	Host string `env:"MONGO_HOST"`
	Port string `env:"MONGO_PORT"`
	Db   string `env:"MONGO_DB"`
	User string `env:"MONGO_USER"`
	Pass string `env:"MONGO_PASS"`
}

func Parse() {
	err := env.Parse(&EveConfig)
	if err != nil {
		panic(err)
	}
	err = env.Parse(&MongoConfig)
	if err != nil {
		panic(err)
	}
}
