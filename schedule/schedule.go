package schedule

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sanyokbig/cats-industry-server/comms"
)

type Schedule struct {
	comms *comms.Comms
	stop  chan bool

	updateJobs int
}

func NewSchedule(comms *comms.Comms, updateJobs int) *Schedule {
	if updateJobs == 0 {
		log.Warning("updateJobs received as 0, setting to 15 instead. Consider setting SCHEDULE_UPDATE_JOBS to non-zero value")
		updateJobs = 15
	}

	return &Schedule{
		comms: comms,
		stop:  make(chan bool),

		updateJobs: updateJobs,
	}
}

func (s *Schedule) Run() {
	updateJobsEvery := time.Minute * time.Duration(s.updateJobs)
	jobsUpdate := time.NewTicker(updateJobsEvery)

	// Initial run
	s.comms.Foreman.UpdateJobs()

	for {
		select {
		case <-jobsUpdate.C:
			s.comms.Foreman.UpdateJobs()
		case <-s.stop:
			log.Info("schedule stop called, exiting")
			break
		}
	}
}
