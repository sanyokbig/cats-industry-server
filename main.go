package main

import (
	"cats-industry-server/server"
	"cats-industry-server/config"
	"cats-industry-server/mongo"
)

func main() {
	config.Parse()
	db := mongo.NewConnection(mongo.GenerateUri(config.MongoConfig.Host, config.MongoConfig.Port, config.MongoConfig.Db, config.MongoConfig.User, config.MongoConfig.Pass))

	err := db.Connect()
	if err != nil {
		panic(err)
	}

	srv := server.Server{
		Db: db,
	}
	srv.Run("9962")
}
