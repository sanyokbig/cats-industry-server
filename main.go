package main

import (
	"cats-industry-server/server"
	"cats-industry-server/config"
	"cats-industry-server/postgres"
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

	srv := server.Server{
		Db: db,
	}
	srv.Run("9962")
}
