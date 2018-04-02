package config

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

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
	Uri        string `env:"REDIS_URI"`
	Pass       string `env:"REDIS_PASS"`
	TTLDays    int    `env:"REDIS_TTL_DAYS"`
	SessionsDB int    `env:"REDIS_DB_SESSIONS"`
	RolesDB    int    `env:"REDIS_DB_ROLES"`
}

var ScheduleConfig struct {
	UpdateJobs int `env:"SCHEDULE_UPDATE_JOBS"`
}

var AppConfig struct {
	Port string `env:"APP_PORT"`
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
	err = env.Parse(&ScheduleConfig)
	if err != nil {
		panic(err)
	}
	err = env.Parse(&AppConfig)
	if err != nil {
		panic(err)
	}
	if AppConfig.Port == "" {
		log.Warningf("app port is not set, consider setting APP_PORT variable")
	}
}
