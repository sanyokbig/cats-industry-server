package methods

import (
	"cats-industry-server/schema"
	"database/sql"
	"log"
)

//easyjson:json
type restoreSessionPayload struct {
	SID string `json:"sid"`
}

// Sends client auth info including user with its characters
func restoreSession(c Client, m schema.Message) (resp *schema.Message, err error) {
	log.Println("restoreSession request from", c.GetID())

	resp = &schema.Message{
		Type:    "restoration",
		Payload: schema.Payload{},
	}

	// Parse incoming payload
	payload := restoreSessionPayload{}
	if err := m.Payload.Deliver(&payload); err != nil {
		log.Println(err)
		return nil, ErrPayloadParseFailed
	}

	// Get userID from session
	userID, err := c.GetComms().Sessions.Get(payload.SID)
	if err != nil {
		return nil, err
	}

	user := &schema.User{}
	err = user.FindWithCharacters(c.GetPostgres(), userID)
	if err != nil && err != sql.ErrNoRows {
		// Unexpected error
		return nil, err
	}
	resp.Payload["user"] = user
	if err == sql.ErrNoRows {
		// Session bound to not existing user, reset session user
		err = c.GetComms().Sessions.Set(payload.SID, 0)
		if err != nil {
			return nil, err
		}

		resp.Payload["user"] = nil
	}

	return resp, nil
}
