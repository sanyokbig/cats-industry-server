package session

import (
	"cats-industry-server/schema"
	"cats-industry-server/server"
)

// Session must be deleted after this time
const SessionLifetime int64 = 86400 * 7 // One week

type Session struct {
	ID      string
	Socket  *server.Client
	User    *schema.User
	Expires int64
}
