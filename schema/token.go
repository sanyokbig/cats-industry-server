package schema

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"cats-industry-server/config"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

//easyjson:json
type Token struct {
	Id           uint `db:"id"`
	UserId       uint
	ExpiresAt    int64  `db:"expires_at"`
	Type         string `json:"type" db:"type"`
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

	log.Println("new", t)
	return nil
}

// Gets token owner
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

	return &o, nil
}

// Updates token from Eve server
func (t Token) IsExpired() bool {
	return time.Now().Unix() >= t.ExpiresAt
}

// Saves token to postgres and updates id in struct
func (t *Token) Save(db *sqlx.DB) error {
	rows, err := db.NamedQuery(`INSERT INTO tokens (expires_at, access_token, refresh_token, type) VALUES (:expires_at, :access_token, :refresh_token, :type)`, t)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}
