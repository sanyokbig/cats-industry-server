package methods

import (
	"cats-industry-server/schema"
)

// Generates new session and returns sid to client
func getSid(c Client, _ schema.Message) (resp *schema.Message, err error) {
	sid, err := c.GetComms().Sessions.New()
	if err != nil {
		return nil, err
	}

	return &schema.Message{
		Type:    "sid",
		Payload: schema.Payload{"sid": sid},
	}, nil
}
