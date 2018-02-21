package auth

import (
	"github.com/sanyokbig/cats-industry-server/schema"

	"github.com/jmoiron/sqlx"
	"github.com/sanyokbig/cats-industry-server/comms"
)

type sessionSender interface {
	SendToSession(session string, message *schema.Message)
}

func notifyClientAboutAuth(db sqlx.Queryer, sid string, userID uint, sender sessionSender, sentinel comms.Sentinel) error {
	// Get full user info
	user := &schema.User{}

	err := user.FindWithCharacters(db, userID)
	if err != nil {
		return err
	}

	roles, err := sentinel.GetRoles(userID)
	if err != nil {
		return err
	}

	// Send auth info to client via comms
	sender.SendToSession(sid, &schema.Message{
		Type: "auth",
		Payload: schema.Payload{
			"user":  user,
			"roles": roles,
		},
	})

	return nil
}
