package server

import (
	"github.com/sanyokbig/cats-industry-server/schema"
	"database/sql"
)

// Sends client auth info including user with its characters
func restoreSession(sid string, hub *Hub) (msg *schema.Message, err error) {
	msg = &schema.Message{
		Type:    "auth",
		Payload: schema.Payload{},
	}

	// Get userID from session
	userID, err := hub.comms.Sessions.Get(sid)
	if err != nil {
		return nil, err
	}

	user := &schema.User{}
	err = user.FindWithCharacters(hub.postgres, userID)
	if err != nil && err != sql.ErrNoRows {
		// Unexpected error
		return nil, err
	}
	msg.Payload["user"] = user
	if err == sql.ErrNoRows {
		// Session bound to not existing user, reset session user
		err = hub.comms.Sessions.Set(sid, 0)
		if err != nil {
			return nil, err
		}

		msg.Payload["user"] = nil
	}

	return msg, nil
}
