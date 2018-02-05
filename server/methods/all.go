package methods

import (
	"cats-industry-server/comms"
	"cats-industry-server/schema"
)

type Handler func(c Client, req schema.Message) (resp *schema.Message, err error)

type Client interface {
	GetID() string
	GetComms() *comms.Comms
}

var all = map[string]Handler{
	"login_request": loginRequest,
	"get_sid":       getSid,
}

func Get(name string) Handler {
	h, ok := all[name]

	if !ok {
		return nil
	}

	return h
}
