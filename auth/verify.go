package auth

import (
	"net/http"
	"encoding/base64"
	"fmt"
	"log"
	"io/ioutil"
	"cats-industry-server/config"
)

func VerifyCode(code string) {
	c := &http.Client{}
	url := fmt.Sprintf("https://login.eveonline.com/oauth/token?grant_type=authorization_code&code=%v", code);
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(config.EveConfig.ClientId+":"+config.EveConfig.SecretKey)))

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
}
