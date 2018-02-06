package methods

import (
	"cats-industry-server/schema"
	"log"
)

//easyjson:json
type restoreSessionPayload struct {
	SID string `json:"sid"`
}

func restoreSession(c Client, m schema.Message) (resp *schema.Message, err error) {
	log.Println("restoreSession request from", c.GetID())

	// Parse incoming payload
	payload := restoreSessionPayload{}
	if err := m.Payload.Deliver(&payload); err != nil {
		log.Println(err)
		return nil, ErrPayloadParseFailed
	}

	userID := c.GetComms().Sessions.Get(payload.SID)

	log.Println(userID)

	return &schema.Message{
		Type: "restoration",
		Payload: schema.Payload{
			"user_id": userID,
		},
	}, nil
}
