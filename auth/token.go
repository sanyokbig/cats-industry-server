package auth

import (
	"fmt"
	"encoding/base64"
	"io/ioutil"
	"encoding/json"
	"net/http"

	"cats-industry-server/config"
	"log"
	"time"
	"github.com/jmoiron/sqlx"
)

type Token struct {
	Id           uint   `db:"id"`
	UserId       uint
	ExpiresAt    int64  `db:"expires_at"`
	Type         string `json:"type" db:"type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token" db:"access_token"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
}

// Creates new token using authorization code
func CreateToken(code string) (token *Token, err error) {
	token = &Token{Id: 1}
	c := &http.Client{}
	url := fmt.Sprintf("https://login.eveonline.com/oauth/token?grant_type=authorization_code&code=%v", code);
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(config.EveConfig.ClientId+":"+config.EveConfig.SecretKey)))

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, token)
	if err != nil {
		return nil, err
	}

	token.ExpiresAt = time.Now().Unix() + 5 //int64(token.ExpiresIn)

	return token, nil
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

	err = json.Unmarshal(bodyBytes, t)
	if err != nil {
		return err
	}

	log.Println("new", t)
	return nil
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
