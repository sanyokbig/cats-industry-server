package server

import (
	"github.com/sanyokbig/cats-industry-server/schema"
	"github.com/sanyokbig/cats-industry-server/server/methods"
	"encoding/json"
	"log"
)

// General processing of ws requests
func processRequest(c *Client, msg []byte) {
	request := schema.Message{}

	err := json.Unmarshal(msg, &request)
	if err != nil {
		log.Println(err)
		return
	}

	if request.Type == "" {
		log.Printf("%v have no \"type\"", string(msg))
		return
	}

	handler := methods.Get(request.Type)
	if handler == nil {
		log.Printf("request \"%v\" not handled: unknown type", request)
		return
	}
	toSend, err := handler(c, request)

	if err != nil {
		log.Println(err)
	}

	if toSend != nil {
		c.Respond(toSend)
	}
}
