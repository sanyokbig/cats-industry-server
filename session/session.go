package session

import (
	"cats-industry-server/comms"
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

type Sessions struct {
	comms *comms.Comms

	// sessionID : userID
	list map[string]uint
}

func New(comms *comms.Comms) *Sessions {
	return &Sessions{
		comms: comms,
		list:  map[string]uint{},
	}
}

func (s *Sessions) Run() {
	for {
		select {}
	}
}
