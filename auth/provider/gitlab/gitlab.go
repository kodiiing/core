package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"kodiiing/auth/provider"
	"net/http"
	"net/url"
)

type gitlab struct {
	ClientId     string
	ClientSecret string
}

func New(id, secret string) provider.Authentication {
	return &gitlab{
		ClientId:     id,
		ClientSecret: secret,
	}
}

type acquireAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int64  `json:"created_at"`
}

func (g *gitlab) AcquireAccessToken(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", provider.ErrCodeEmpty
	}

	requestQuery := url.Values{}
	requestQuery.Set("client_id", g.ClientId)
	requestQuery.Set("client_secret", g.ClientSecret)
	requestQuery.Set("code", code)
	requestQuery.Set("grant_type", "authorization_code")

	requestUrl := url.URL{
		Scheme:   "https",
		Host:     "gitlab.com",
		Path:     "/oauth/token",
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
