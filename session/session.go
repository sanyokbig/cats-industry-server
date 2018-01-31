package session

import (
	"cats-industry-server/server/ws"
	"cats-industry-server/schema"
)

type Session struct {
	socket *ws.Client
	user   *schema.User
}
