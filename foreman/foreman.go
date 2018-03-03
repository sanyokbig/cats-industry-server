package foreman

import (
	"log"

	"github.com/sanyokbig/cats-industry-server/comms"
)

type Foreman struct {
	comms *comms.Comms
}

func NewForeman(comms *comms.Comms) *Foreman {
	return &Foreman{
		comms: comms,
	}
}

func (f Foreman) UpdateJobs() {
	log.Println("updateing jobs")
}
