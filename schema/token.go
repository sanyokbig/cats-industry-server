package schema

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"time"

	"github.com/sanyokbig/cats-industry-server/config"
	log "github.com/sirupsen/logrus"

	"github.com/sanyokbig/cats-industry-server/postgres"

	"github.com/go-errors/errors"
)

//easyjson:json
type Token struct {
	ID           uint   `db:"id"`
	CharacterID  uint   `db:"character_id"`
	ExpiresAt    int64  `db:"expires_at"`
	Type         string `db:"type"`
	Scopes       string `db:"scopes"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token" db:"access_token"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
}

//easyjson:json
type Owner struct {
	CharacterID        uint   `json:"CharacterID"`
	CharacterName      string `json:"CharacterName"`
	ExpiresOn          string `json:"ExpiresOn"`
	Scopes             string `json:"Scopes"`
	TokenType          string `json:"TokenType"`
	CharacterOwnerHash string `json:"CharacterOwnerHash"`
}

// Updates token from Eve server
func (t *Token) Refresh() error {
	c := &http.Client{}
	url := fmt.Sprintf("https://login.eveonline.com/oauth/token?grant_type=refresh_token&refresh_token=%v", t.RefreshToken)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(config.EveConfig.ClientId+":"+config.EveConfig.SecretKey)))

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = t.UnmarshalJSON(bodyBytes)
	if err != nil {
		return err
	}

	log.Debugf("refreshed token %v", t)
	return nil
}

// Gets token owner, updates token information and returns Owner
func (t *Token) GetOwner() (*Owner, error) {
	c := &http.Client{}
	url := fmt.Sprintf("https://login.eveonline.com/oauth/verify")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+t.AccessToken)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	o := Owner{}

	err = o.UnmarshalJSON(bodyBytes)
	if err != nil {
		return nil, err
	}

	// Update token fields from received owner
	t.CharacterID = o.CharacterID
	t.Scopes = o.Scopes

	return &o, nil
}

// Updates token from Eve server
func (t Token) IsExpired() bool {
	return time.Now().Unix() >= t.ExpiresAt
}

// Saves token to postgres and updates id in struct
func (t *Token) Save(db postgres.NamedQueryer) error {
	rows, err := db.NamedQuery(`
		INSERT INTO tokens 
			(character_id, expires_at, access_token, refresh_token, scopes) 
		VALUES 
			(:character_id, :expires_at, :access_token, :refresh_token, :scopes) 
		ON CONFLICT (character_id, scopes) DO UPDATE 
			SET (expires_at, access_token, refresh_token) = (:expires_at, :access_token, :refresh_token) 
		RETURNING id`,
		t,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return errors.New("token create: rows.Scan() failed")
	}
	err = rows.StructScan(t)
	if err != nil {
		return err
	}

	return nil
}
