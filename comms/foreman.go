package comms

import "github.com/sanyokbig/cats-industry-server/schema"

type Foreman interface {
	UpdateJobs()
	GetJobs() (*schema.Jobs, error)
}
