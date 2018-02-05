package methods

import (
	"cats-industry-server/config"
	"cats-industry-server/schema"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/satori/go.uuid"
	"log"
)

var scopeSets = map[string]string{
	"simple":        "publicData",
	"industrialist": "publicData characterIndustryJobsRead corporationIndustryJobsRead esi-mail.send_mail.v1",
}

var (
	ErrPayloadParseFailed = errors.New("failed to parse payload")
	ErrNoScopeSet         = errors.New("scope set not found")
	ErrStateGenFailed     = errors.New("failed to generate state")
)

//easyjson:json
type loginRequestPayload struct {
	ScopeSet string `json:"scope_set"`
}

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
		log.Println(err)
		return nil, ErrNoScopeSet
	}

	// Generate state to identify login response later
	state, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		return nil, ErrStateGenFailed
	}

	// Generate login uri for client
	uri := fmt.Sprintf("https://login.eveonline.com/oauth/authorize?response_type=code&redirect_uri=%v&client_id=%v&scope=%v&state=%v", config.EveConfig.RedirectUri, config.EveConfig.ClientId, scopes, state)

	return &schema.Message{
		Type:    "login_uri",
		Payload: schema.Payload{"uri": uri},
	}, nil
}
