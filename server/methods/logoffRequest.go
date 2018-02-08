package methods

import (
	"cats-industry-server/schema"
	"log"
)

// Generate login uri for client and add client to pending
func logoffRequest(c Client, m schema.Message) (resp *schema.Message, err error) {
	log.Println("logoff request from", c.GetID())

	resp = &schema.Message{}

	err = c.GetComms().Sessions.Set(c.GetSID(), 0)
	if err != nil {
		resp.Type = "logoff_fail"
		return resp, err
	}

	resp.Type = "logoff_ok"
	return resp, nil
}
