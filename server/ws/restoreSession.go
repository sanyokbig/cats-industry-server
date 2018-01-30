package ws

import (
	"log"
)

func restoreSession(responder Responder, req Message) (resp *Message, err error) {
	log.Println("log request from", responder)
	return &Message{
		Type:    "login_uri",
		Payload: Payload{"uri": "http://google.com"},
	}, nil
}
