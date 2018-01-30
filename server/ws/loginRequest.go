package ws

import (
	"log"
	"fmt"
	"cats-industry-server/config"
)

func loginRequest(c *Client, _ Message) (resp *Message, err error) {
	log.Println("log request from", c.id)

	uri := fmt.Sprintf("https://login.eveonline.com/oauth/authorize?response_type=code&redirect_uri=%v&client_id=%v&scope=%v&state=%v", config.EveConfig.RedirectUri, config.EveConfig.ClientId, "publicData", "so-state-very-unique")

	return &Message{
		Type:    "login_uri",
		Payload: Payload{"uri": uri},
	}, nil
}
