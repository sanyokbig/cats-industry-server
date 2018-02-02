package session

import (
	"cats-industry-server/server"
	"cats-industry-server/schema"
)

type Session struct {
	socket *server.Client
	user   *schema.User
}
