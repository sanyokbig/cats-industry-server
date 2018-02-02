package methods

import (
	"github.com/satori/go.uuid"
	"cats-industry-server/schema"
)

func getSid(c Client, req schema.Message) (resp *schema.Message, err error) {
	return &schema.Message{
		Type:    "sid",
		Payload: schema.Payload{"sid": uuid.Must(uuid.NewV1()).String()},
	}, nil
}
