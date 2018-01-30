package ws

type Handler func(c *Client, req Message) (resp *Message, err error)

type Responder interface {
	Respond(msg Message)
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
