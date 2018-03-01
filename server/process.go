package server

import (
	"encoding/json"

	"github.com/sanyokbig/cats-industry-server/schema"
	"github.com/sanyokbig/cats-industry-server/server/methods"
	log "github.com/sirupsen/logrus"
)

// General processing of ws requests
func processRequest(c *Client, msg []byte) {
	request := schema.Message{}

	err := json.Unmarshal(msg, &request)
	if err != nil {
		log.Errorf("failed to unmarshal request: %v", err)
		return
	}

	if request.Type == "" {
		log.Errorf("%v have no 'type'", string(msg))
		return
	}

	handler := methods.Get(request.Type)
	if handler == nil {
		log.Warningf("request '%v' not handled: unknown type %v", request, request.Type)
		return
	}

	toSend, err := handler(c, request)
	if err != nil {
		log.Errorf("failed while executing handler: %v", err)
	}

	if toSend != nil {
		c.Respond(toSend)
	}
}
