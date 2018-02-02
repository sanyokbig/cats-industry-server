package methods

import (
	"log"
	"fmt"
	"cats-industry-server/config"
	"cats-industry-server/schema"
)

func loginRequest(c Client, _ schema.Message) (resp *schema.Message, err error) {
	log.Println("log request from", c.GetID())

	uri := fmt.Sprintf("https://login.eveonline.com/oauth/authorize?response_type=code&redirect_uri=%v&client_id=%v&scope=%v&state=%v", config.EveConfig.RedirectUri, config.EveConfig.ClientId, "publicData", "so-state-very-unique")

	return &schema.Message{
		Type:    "login_uri",
		Payload: schema.Payload{"uri": uri},
	}, nil
}
