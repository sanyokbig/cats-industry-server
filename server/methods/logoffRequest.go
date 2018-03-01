package methods

import (
	"github.com/sanyokbig/cats-industry-server/schema"
	log "github.com/sirupsen/logrus"
)

// Generate login uri for client and add client to pending
func logoffRequest(c Client, _ schema.Message) (resp *schema.Message, err error) {
	log.Infof("logoff request from %v", c.GetID())

	resp = &schema.Message{}

	err = c.GetComms().Sessions.Set(c.GetSID(), 0)
	if err != nil {
		resp.Type = "logoff_fail"
		return resp, err
	}

	resp.Type = "logoff_ok"
	resp.Payload.SetAsDefaultAuthPayload()
	return resp, nil
}
