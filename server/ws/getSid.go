package ws

import (
	"log"
	"github.com/satori/go.uuid"
)

func getSid(c *Client, req Message) (resp *Message, err error) {
	log.Println("get_sid request from", c.id)

	return &Message{
		Type:    "sid",
		Payload: Payload{"sid": uuid.Must(uuid.NewV1()).String()},
	}, nil
}
