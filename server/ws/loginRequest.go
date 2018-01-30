package ws

import (
	"log"
	"fmt"
	"github.com/caarlos0/env"
)

type EveConfig struct {
	ClientId    string `env:"CLIENT_ID"`
	SecretKey   string `env:"SECRET_KEY"`
	RedirectUri string `env:"REDIRECT_URI"`
}

func loginRequest(c *Client, _ Message) (resp *Message, err error) {
	log.Println("log request from", c.id)

	conf := EveConfig{}

	err = env.Parse(&conf)
	if err!= nil {
		panic(err)
	}

	log.Println(conf)

	uri := fmt.Sprintf("https://login.eveonline.com/oauth/authorize?response_type=code&redirect_uri=%v&client_id=%v&scope=%v&state=%v", conf.RedirectUri, conf.ClientId, "publicData","so-state-very-unique")

	return &Message{
		Type:    "login_uri",
		Payload: Payload{"uri": uri},
	}, nil
}
