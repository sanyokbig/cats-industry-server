package server

import (
	"github.com/sanyokbig/cats-industry-server/schema"
	log "github.com/sirupsen/logrus"
)

func prepareSession(hub *Hub) (msg *schema.Message, sid string, err error) {
	sid, err = hub.comms.Sessions.New()
	if err != nil {
		log.Errorf("failed to create new session: %v", err)
		return
	}

	msg = &schema.Message{
		Type:    "sid",
		Payload: schema.Payload{"sid": sid},
	}

	return
}
