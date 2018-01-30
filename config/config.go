package config

import "github.com/caarlos0/env"

var EveConfig struct {
	ClientId    string `env:"CLIENT_ID"`
	SecretKey   string `env:"SECRET_KEY"`
	RedirectUri string `env:"REDIRECT_URI"`
}

func Parse() {
	err := env.Parse(&EveConfig)
	if err != nil {
		panic(err)
	}
}
