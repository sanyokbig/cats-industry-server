package foreman

import (
	log "github.com/sirupsen/logrus"

	"sync"

	"fmt"

	"net/http"

	"io/ioutil"

	"github.com/go-errors/errors"
	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/postgres"
	"github.com/sanyokbig/cats-industry-server/schema"
)

type Foreman struct {
	comms *comms.Comms

	db      postgres.DB
	baseUri string
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
	jobs, err := f.pullJobs(tokens)
	if err != nil {
		log.Errorf("failed to pull jobs: %v", err)
		return
	}

	log.Debugf("pulled jobs: %v", len(*jobs))

	// Save to db
	err = jobs.Save(f.db)
	if err != nil {
		log.Errorf("failed to save jobs: %v", err)
	}
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
func (f *Foreman) pullJobs(tokens *schema.Tokens) (jobs *schema.Jobs, err error) {
	wg := sync.WaitGroup{}
	pulledJobs := make(chan *schema.Jobs, len(*tokens))

	for _, t := range *tokens {
		wg.Add(1)
		go f.goPull(&wg, t, pulledJobs)
	}

	wg.Wait()
	close(pulledJobs)

	jobs = &schema.Jobs{}
	for js := range pulledJobs {
		if js != nil {
			*jobs = append(*jobs, *js...)
		}
	}

	return jobs, nil
}

func (f *Foreman) goPull(wg *sync.WaitGroup, t schema.Token, result chan<- *schema.Jobs) {
	defer wg.Done()
	jobs, err := f.useToken(&t)
	if err != nil {
		log.Errorf("failed to pull jobs with token %v: %v", t.ID, err)
	}
	result <- jobs
}

// Pull jobs with passed token
func (f *Foreman) useToken(token *schema.Token) (jobs *schema.Jobs, err error) {
	// Make sure token is alive
	if err := token.Refresh(f.db); err != nil {
		log.Error("failed to refresh token: %v", err)
		return nil, schema.ErrFailedToRefreshToken
	}

	uri := fmt.Sprintf("https://esi.tech.ccp.is/latest/characters/%v/industry/jobs/?include_completed=true&token=%v",
		token.CharacterID, token.AccessToken,
	)

	resp, err := http.Get(uri)
	if err != nil {
		log.Errorf("failed to get: %v", err)
		return nil, errors.New("failed to get jobs")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read body: %v", err)
		return nil, errors.New("failed to read body")
	}

	jobs = &schema.Jobs{}
	err = jobs.UnmarshalJSON(body)
	if err != nil {
		log.Errorf("failed to unmarshal body: %v, body: %s", err, body)
		return nil, errors.New("failed to unmarshal body")
	}

	return jobs, nil
}
