package foreman

import (
	log "github.com/sirupsen/logrus"

	"sync"

	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/postgres"
	"github.com/sanyokbig/cats-industry-server/schema"
)

type Foreman struct {
	comms *comms.Comms

	db postgres.DB
}

func NewForeman(comms *comms.Comms, db postgres.DB) *Foreman {
	return &Foreman{
		comms: comms,
		db:    db,
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

	jobs, err := f.useTokens(tokens)
	if err != nil {
		log.Errorf("failed to pull jobs: %v", err)
		return
	}
	log.Debugf("pulled jobs: %v", jobs)

	// Upsert to db
}

// Get all industrial tokens from db
func (f *Foreman) getTokens() (*schema.Tokens, error) {
	tokens := &schema.Tokens{}
	err := tokens.GetTokensOfScope(f.db, "industrial")
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Use passed tokens to get jobs list from EVE server
func (f *Foreman) useTokens(tokens *schema.Tokens) (jobs *schema.Jobs, err error) {
	wg := sync.WaitGroup{}
	pulledJobs := make(chan *schema.Jobs, len(*tokens))

	for _, t := range *tokens {
		wg.Add(1)
		go f.goPull(&wg, t, pulledJobs)
	}

	wg.Wait()
	close(pulledJobs)

	//jobs := &schema.Jobs{}
	for js := range pulledJobs {
		log.Info(js)
	}

	return jobs, nil
}

func (f *Foreman) goPull(wg *sync.WaitGroup, t schema.Token, result chan<- *schema.Jobs) {
	defer wg.Done()
	jobs, err := f.pullJobs(&t)
	if err != nil {
		log.Errorf("failed to pull jobs with token %v: %v", t.ID, err)
	}
	result <- jobs
}

// Pull jobs with passed token
func (f *Foreman) pullJobs(token *schema.Token) (jobs *schema.Jobs, err error) {
	// Make sure token is alive
	if err := token.Refresh(f.db); err != nil {
		log.Error(err)
		return nil, schema.ErrFailedToRefreshToken
	}

	return jobs, nil
}
