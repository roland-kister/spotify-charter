package spotify

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
)

const authURL = "https://accounts.spotify.com/api/token"

type AuthInterceptor struct {
	core        http.RoundTripper
	accessToken *string
}

type AuthResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"toke_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (a AuthInterceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("Authorization", "Bearer "+*a.accessToken)

	return a.core.RoundTrip(r)
}

func (c *APICLient) Authorize() error {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+basicAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.DefaultClient

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return authErrRespToErr(&res.Body)
	}

	authResp, err := decodeResp[AuthResp](&res.Body)
	if err != nil {
		return err
	}

	c.accessToken = authResp.AccessToken

	return nil
}
