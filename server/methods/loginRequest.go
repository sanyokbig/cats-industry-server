package methods

import (
	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/config"
	"github.com/sanyokbig/cats-industry-server/schema"
	"fmt"
	"log"

	"github.com/go-errors/errors"
)

var scopeSets = map[string]string{
	"simple":     "publicData",
	"industrial": "characterIndustryJobsRead corporationIndustryJobsRead",
	"mailing":    "esi-mail.send_mail.v1",
}

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
	log.Println("log request from", c.GetID())

	// Parse incoming payload
	payload := loginRequestPayload{}
	if err := m.Payload.Deliver(&payload); err != nil {
		log.Println(err)
		return nil, ErrPayloadParseFailed
	}

	// Select scope set for authorization
	scopes, ok := scopeSets[payload.ScopeSet]
	if !ok {
		log.Printf("scope set for %v not found", payload.ScopeSet)
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
