package main

import (
	"cats-industry-server/server"
	"cats-industry-server/config"
)

func main() {
	config.Parse()
	server.Run("9962")
}
