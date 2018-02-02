package methods

import (
	"log"
	"cats-industry-server/schema"
)

func restoreSession(c Client, req schema.Message) (resp *schema.Message, err error) {
	log.Println("log request from", c.GetID())
	return &schema.Message{
		Type:    "login_uri",
		Payload: schema.Payload{"uri": "http://google.com"},
	}, nil
}
