package ws

import (
	"encoding/json"
	"log"
)

// Looks for appropriate handler and calls it
func processRequest(c *Client, msg []byte) {
	request := Message{}

	err := json.Unmarshal(msg, &request)
	if err != nil {
		log.Println(err)
		return
	}

	if request.Type == "" {
		log.Printf("%v have no \"type\"", string(msg))
		return
	}

	handler := Get(request.Type)

	if handler == nil {
		log.Printf("request \"%v\" not handled: unknown type", request)
		return
	}

	toSend, err := handler(c, request)

	if toSend != nil {
		c.Respond(*toSend)
	}
}
