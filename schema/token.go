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
	"github.com/jmoiron/sqlx"
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

type Tokens []Token

//easyjson:json
type Owner struct {
	CharacterID        uint   `json:"CharacterID"`
	CharacterName      string `json:"CharacterName"`
	ExpiresOn          string `json:"ExpiresOn"`
	Scopes             string `json:"Scopes"`
	TokenType          string `json:"TokenType"`
	CharacterOwnerHash string `json:"CharacterOwnerHash"`
}

var (
	ErrFailedToRefreshToken = errors.New("failed to refresh token")
)

// Refreshes token if needed
func (t *Token) Refresh(db postgres.NamedQueryer) error {
	if !t.IsExpired() {
		return nil
	}
	log.Debugf("token %v expired, refreshing", t.ID)
	err := t.refresh()
	if err != nil {
		return nil
	}

	t.ExpiresAt = time.Now().Unix() + int64(t.ExpiresIn)

	err = t.Save(db)
	if err != nil {
		log.Warningf("failed to save refreshed token: %v", err)
	}

	return nil
}

// Updates token from Eve server
func (t *Token) refresh() error {
	c := &http.Client{}
	url := fmt.Sprintf("https://login.eveonline.com/oauth/token?grant_type=refresh_token&refresh_token=%v", t.RefreshToken)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error(err)
		return errors.New("failed to prepare post request")
	}

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(config.EveConfig.ClientId+":"+config.EveConfig.SecretKey)))

	resp, err := c.Do(req)
	if err != nil {
		log.Error(err)
		return errors.New("failed to do post request")
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return errors.New("failed to read response")
	}

	err = t.UnmarshalJSON(bodyBytes)
	if err != nil {
		log.Error(err)
		return errors.New("failed to unmarshal response")
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

func (t *Tokens) GetTokensOfScope(queryer sqlx.Queryer, setName string) error {
	set, ok := ScopeSets[setName]
	if !ok {
		return errors.New(fmt.Sprintf("scope set with name '%v' not found", setName))
	}
	rows, err := queryer.Queryx(`
		SELECT id, character_id, expires_at, scopes, access_token, refresh_token FROM tokens WHERE scopes = $1;
	`, set)
	if err != nil {
		return err
	}
	defer rows.Close()

	token := Token{}
	for rows.Next() {
		err := rows.StructScan(&token)
		if err != nil {
			return err
		}
		*t = append(*t, token)
	}

	return nil
}
