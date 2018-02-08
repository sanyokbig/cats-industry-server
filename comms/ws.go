package comms

import "cats-industry-server/schema"

type Hub interface {
	SendToSession(session string, msg *schema.Message)
}
