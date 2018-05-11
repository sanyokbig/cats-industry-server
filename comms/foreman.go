package comms

import "github.com/sanyokbig/cats-industry-server/schema"

type Foreman interface {
	UpdateJobs()
	GetJobs(params schema.GetParams) (*schema.Jobs, error)
}
