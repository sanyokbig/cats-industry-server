package foreman

import (
	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/schema"
)

type Foreman struct {
	comms *comms.Comms

	db sqlx.Queryer
}

func NewForeman(comms *comms.Comms, queryer sqlx.Queryer) *Foreman {
	return &Foreman{
		comms: comms,
		db:    queryer,
	}
}

func (f *Foreman) UpdateJobs() {
	log.Infof("updating jobs")

	tokens, err := f.getTokens()
	if err != nil {
		log.Errorf("failed to get tokens: %v", err)
		return
	}

	log.Debugf("industrial tokens: %v", tokens)
	// Get all jobs using token

	jobs, err := f.pullJobs(tokens)
	if err != nil {
		log.Errorf("failed to pull jobs: %v", err)
		return
	}
	log.Debugf("pulled jobs: %v", jobs)

	// Upsert to db
}

func (f *Foreman) getTokens() (*schema.Tokens, error) {
	tokens := &schema.Tokens{}
	err := tokens.GetTokensOfScope(f.db, "industrial")
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (f *Foreman) pullJobs(tokens *schema.Tokens) (*schema.Jobs, error) {
	return nil, nil
}
