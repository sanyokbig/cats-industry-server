package methods

import (
	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/postgres"
	"github.com/sanyokbig/cats-industry-server/schema"
)

type Handler func(c Client, req schema.Message) (resp *schema.Message, err error)

type Client interface {
	GetID() string
	GetSID() string
	GetComms() *comms.Comms
	GetPostgres() *postgres.Connection
}

var all = map[string]Handler{
	"login_request":  loginRequest,
	"logoff_request": logoffRequest,
	"get_jobs":       getJobs,
}

func Get(name string) Handler {
	h, ok := all[name]

	if !ok {
		return nil
	}

	return h
}
