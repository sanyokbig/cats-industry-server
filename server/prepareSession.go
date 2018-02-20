package server

import (
	"log"
	"github.com/sanyokbig/cats-industry-server/schema"
)

func prepareSession(hub *Hub) (msg *schema.Message, err error) {
	sid, err := hub.comms.Sessions.New()
	if err != nil {
		log.Print(err)
		return
	}

	msg = &schema.Message{
		Type:    "sid",
		Payload: schema.Payload{"sid": sid},
	}

	return msg, err
}
