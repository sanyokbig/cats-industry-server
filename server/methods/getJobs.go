package methods

import (
	"github.com/sanyokbig/cats-industry-server/schema"
	log "github.com/sirupsen/logrus"
)

func getJobs(c Client, req schema.Message) (resp *schema.Message, err error) {
	log.Infof("get jobs request from session %v", c.GetSID())

	resp = &schema.Message{}

	allowed := c.GetComms().Sentinel.CheckSession(c.GetSID(), "see_shared_jobs")
	if !allowed {
		resp.Type = "get_jobs_denied"
		return resp, nil
	}

	jobs, err := c.GetComms().Foreman.GetJobs()
	if err != nil {
		resp.Type = "get_jobs_fail"
		return resp, err
	}

	err = resp.Payload.Pack(jobs)
	if err != nil {
		resp.Type = "get_jobs_fail"
		return resp, err
	}

	resp.Type = "get_jobs_ok"

	log.Infof("get jobs success for session %v", c.GetSID())

	return resp, nil
}
