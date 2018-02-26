package server

import (
	"github.com/sanyokbig/cats-industry-server/schema"
	"database/sql"
)

// Sends client auth info including user with its characters
func restoreSession(sid string, hub *Hub) (msg *schema.Message, err error) {
	msg = &schema.Message{
		Type: "auth",
	}

	// Get userID from session
	userID, err := hub.comms.Sessions.Get(sid)
	if err != nil {
		return nil, err
	}

	// Get auth payload for user if found
	payload, err := schema.User{ID: userID}.GetAuthPayload(hub.postgres, hub.comms.Sentinel)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	msg.Payload = *payload

	// Check if session bound to not existing user
	if err == sql.ErrNoRows {
		// Reset session user
		err = hub.comms.Sessions.Set(sid, 0)
		if err != nil {
			return nil, err
		}
		msg.Payload.SetAsDefaultAuthPayload()
	}

	return msg, nil
}
