package comms

import "github.com/sanyokbig/cats-industry-server/schema"

type Hub interface {
	SendToSession(session string, msg *schema.Message)
}
