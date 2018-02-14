package main

import (
	"github.com/sanyokbig/cats-industry-server/config"
	"github.com/sanyokbig/cats-industry-server/postgres"
	"github.com/sanyokbig/cats-industry-server/server"

	"github.com/go-redis/redis"
)

func main() {
	config.Parse()
	db := postgres.NewConnection(
		config.PostgresConfig.Host,
		config.PostgresConfig.Port,
		config.PostgresConfig.Db,
		config.PostgresConfig.User,
		config.PostgresConfig.Pass,
	)

	err := db.Connect()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:       config.RedisConfig.Uri,
		DB:         config.RedisConfig.DB,
		Password:   config.RedisConfig.Pass,
		MaxRetries: 5,
	})
	_, err = client.Ping().Result()
	if err != nil {
		panic(err)
	}

	srv := server.Server{
		Postgres: db,
		Redis:    client,
	}

	srv.Run("9962")
}
