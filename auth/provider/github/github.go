package github

import (
	"context"
	"encoding/json"
	"fmt"
	"kodiiing/auth/provider"
	"net/http"
	"net/url"
)

type github struct {
	ClientId     string
	ClientSecret string
}

func New(id, secret string) provider.Authentication {
	return &github{
		ClientId:     id,
		ClientSecret: secret,
	}
}

type acquireAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (g *github) AcquireAccessToken(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", provider.ErrCodeEmpty
	}

	requestQuery := url.Values{}
	requestQuery.Set("client_id", g.ClientId)
	requestQuery.Set("client_secret", g.ClientSecret)
	requestQuery.Set("code", code)

	requestUrl := url.URL{
		Scheme:   "https",
		Host:     "github.com",
		Path:     "/login/oauth/access_token",
		RawQuery: requestQuery.Encode(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error response: %d", resp.StatusCode)
	}

	var response acquireAccessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	return response.AccessToken, nil
}
