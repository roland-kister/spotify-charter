package spotify

import (
	"encoding/json"
	"io"
	"net/http"
)

const baseURL = "https://api.spotify.com"

type ApiClient struct {
	httpClient   *http.Client
	clientID     string
	clientSecret string
	accessToken  string
}

func NewApiClient(clientID string, clientSecret string) *ApiClient {
	client := &ApiClient{
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	client.httpClient = &http.Client{
		Transport: AuthInterceptor{
			core:        http.DefaultTransport,
			accessToken: &client.accessToken,
		},
	}

	return client
}

func decodeResp[T interface{}](body *io.ReadCloser) (*T, error) {
	var resp T

	if err := json.NewDecoder(*body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
