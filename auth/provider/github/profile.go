package github

import (
	"context"
	"encoding/json"
	"fmt"
	"kodiiing/auth"
	"kodiiing/auth/provider"
	"net/http"
	"net/url"
	"time"
)

type getProfileResponse struct {
	Login       string `json:"login"`
	ID          int64  `json:"id"`
	NodeID      string `json:"node_id"`
	AvatarUrl   string `json:"avatar_url"`
	GravatarID  string `json:"gravatar_id"`
	URL         string `json:"url"`
	HtmlURL     string `json:"html_url"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Email       string `json:"email"`
	PublicRepos int64  `json:"public_repos"`
	Followers   int64  `json:"followers"`
	Following   int64  `json:"following"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (g *github) GetProfile(ctx context.Context, accessToken string) (auth.User, error) {
	if accessToken == "" {
		return auth.User{}, provider.ErrCodeEmpty
	}

	// For reference, see:
	// https://docs.github.com/en/rest/users/users#get-the-authenticated-user
	requestUrl := url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "/user",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return auth.User{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return auth.User{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return auth.User{}, fmt.Errorf("error response: %d", resp.StatusCode)
	}

	var responseBody getProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return auth.User{}, fmt.Errorf("error decoding response: %w", err)
	}

	userAvatarUrl, err := url.Parse(responseBody.AvatarUrl)
	if err != nil {
		return auth.User{}, fmt.Errorf("error parsing avatar url: %w", err)
	}

	userProfileUrl, err := url.Parse(responseBody.HtmlURL)
	if err != nil {
		return auth.User{}, fmt.Errorf("error parsing profile url: %w", err)
	}

	userCreatedAt, err := time.Parse(time.RFC3339, responseBody.CreatedAt)
	if err != nil {
		return auth.User{}, fmt.Errorf("error parsing created at: %w", err)
	}

	return auth.User{
		ID:               responseBody.ID,
		Provider:         auth.ProviderGithub,
		NodeID:           responseBody.NodeID,
		Name:             responseBody.Name,
		Username:         responseBody.Login,
		AvatarURL:        userAvatarUrl,
		ProfileURL:       userProfileUrl,
		Location:         responseBody.Location,
		Email:            responseBody.Email,
		PublicRepository: responseBody.PublicRepos,
		Followers:        responseBody.Followers,
		Following:        responseBody.Following,
		CreatedAt:        userCreatedAt,
		RegisteredAt:     time.Time{},
	}, nil
}
