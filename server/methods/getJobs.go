package methods

import (
	"github.com/sanyokbig/cats-industry-server/schema"
	log "github.com/sirupsen/logrus"
)

func getJobs(c Client, req schema.Message) (resp *schema.Message, err error) {
	log.Infof("get jobs request from session %v", c.GetSID())
	resp = schema.NewMessage()

	allowed := c.GetComms().Sentinel.CheckSession(c.GetSID(), "see_shared_jobs")
	if !allowed {
		resp.Type = "get_jobs_denied"
		return resp, nil
	}

	getParams := schema.GetParams{}
	err = req.Payload.Deliver(&getParams)
	if err != nil {
		resp.Type = "get_jobs_fail"
		return resp, err
	}

	jobs, err := c.GetComms().Foreman.GetJobs(getParams)
	if err != nil {
		resp.Type = "get_jobs_fail"
		return resp, err
	}

	resp.Payload["jobs"] = jobs
	resp.Type = "get_jobs_ok"

	log.Infof("get jobs success for session %v", c.GetSID())

	return resp, nil
}
