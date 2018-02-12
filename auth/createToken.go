package auth

import (
	"cats-industry-server/config"
	"cats-industry-server/schema"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Creates new token using authorization code
func CreateToken(code string) (token *schema.Token, err error) {
	token = &schema.Token{}
	c := &http.Client{}
	url := fmt.Sprintf("https://login.eveonline.com/oauth/token?grant_type=authorization_code&code=%v", code)
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

	err = token.UnmarshalJSON(bodyBytes)
	if err != nil {
		return nil, err
	}

	token.ExpiresAt = time.Now().Unix() + 5 //int64(token.ExpiresIn)

	return token, nil
}
