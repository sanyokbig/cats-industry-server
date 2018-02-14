package auth

import (
	"github.com/sanyokbig/cats-industry-server/schema"

	"github.com/jmoiron/sqlx"
)

type sessionSender interface {
	SendToSession(session string, message *schema.Message)
}

func notifyClientAboutAuth(db sqlx.Queryer, sid string, userID uint, sender sessionSender) error {
	// Get full user info
	user := &schema.User{}

	err := user.FindWithCharacters(db, userID)
	if err != nil {
		return err
	}

	// Send auth info to client via comms
	sender.SendToSession(sid, &schema.Message{
		Type: "auth",
		Payload: schema.Payload{
			"user": user,
		},
	})

	return nil
}
