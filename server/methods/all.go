package methods

import (
	"cats-industry-server/comms"
	"cats-industry-server/postgres"
	"cats-industry-server/schema"
)

type Handler func(c Client, req schema.Message) (resp *schema.Message, err error)

type Client interface {
	GetID() string
	GetComms() *comms.Comms
	GetPostgres() *postgres.Connection
}

var all = map[string]Handler{
	"login_request":   loginRequest,
	"get_sid":         getSid,
	"restore_session": restoreSession,
}

func Get(name string) Handler {
	h, ok := all[name]

	if !ok {
		return nil
	}

	return h
}
