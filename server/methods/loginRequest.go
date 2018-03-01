package methods

import (
	"fmt"

	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/config"
	"github.com/sanyokbig/cats-industry-server/schema"
	log "github.com/sirupsen/logrus"

	"github.com/go-errors/errors"
)

var (
	ErrPayloadParseFailed = errors.New("failed to parse payload")
	ErrNoScopeSet         = errors.New("scope set not found")
	ErrStateGenFailed     = errors.New("failed to generate state")
)

//easyjson:json
type loginRequestPayload struct {
	ScopeSet string `json:"scope_set"`
	SID      string `json:"sid"`
}

// Generate login uri for client and add client to pending
func loginRequest(c Client, m schema.Message) (resp *schema.Message, err error) {
	log.Infof("log request from %v", c.GetID())

	// Parse incoming payload
	payload := loginRequestPayload{}
	if err := m.Payload.Deliver(&payload); err != nil {
		log.Errorf("failed to deliver payload: %v", err)
		return nil, ErrPayloadParseFailed
	}

	// Select scope set for authorization
	scopes, ok := schema.ScopeSets[payload.ScopeSet]
	if !ok {
		log.Errorf("no scope set found for '%v'", payload.ScopeSet)
		return nil, ErrNoScopeSet
	}

	// Set state to client id identify login response later
	state := payload.SID

	// Add state to pending
	c.GetComms().Pending.Add <- comms.PendingAdd{State: state, Client: c.GetID()}

	// Generate login uri for client
	uri := fmt.Sprintf("https://login.eveonline.com/oauth/authorize?response_type=code&redirect_uri=%v&client_id=%v&scope=%v&state=%v", config.EveConfig.RedirectUri, config.EveConfig.ClientId, scopes, state)

	return &schema.Message{
		Type:    "login_uri",
		Payload: schema.Payload{"uri": uri},
	}, nil
}
