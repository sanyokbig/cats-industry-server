package main

import (
	"github.com/go-redis/redis"
	"github.com/sanyokbig/cats-industry-server/config"
	"github.com/sanyokbig/cats-industry-server/postgres"
	"github.com/sanyokbig/cats-industry-server/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetLevel(log.DebugLevel)

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
		log.Fatalln("failed to connect with pg:", err)
	}

	redisClients, err := getRedisClients()
	if err != nil {
		log.Fatalln("failed to get redis clients:", err)
	}

	srv := server.Server{
		Postgres:     db,
		RedisClients: redisClients,
	}

	srv.Run("9962")
}

func getRedisClients() (*server.RedisClients, error) {
	clients := &server.RedisClients{}

	// Sessions
	sessions := redis.NewClient(&redis.Options{
		Addr:       config.RedisConfig.Uri,
		DB:         config.RedisConfig.SessionsDB,
		Password:   config.RedisConfig.Pass,
		MaxRetries: 5,
	})
	_, err := sessions.Ping().Result()
	if err != nil {
		return nil, err
	}
	clients.Sessions = sessions

	// Roles
	roles := redis.NewClient(&redis.Options{
		Addr:       config.RedisConfig.Uri,
		DB:         config.RedisConfig.RolesDB,
		Password:   config.RedisConfig.Pass,
		MaxRetries: 5,
	})
	_, err = sessions.Ping().Result()
	if err != nil {
		return nil, err
	}
	clients.Roles = roles

	return clients, nil
}
