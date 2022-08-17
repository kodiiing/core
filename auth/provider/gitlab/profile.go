package gitlab

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
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	WebUrl    string `json:"web_url"`
	CreatedAt string `json:"created_at"`
	Location  string `json:"location"`
	Followers int64  `json:"followers"`
	Following int64  `json:"following"`
	Email     string `json:"email"`
}

func (g *gitlab) GetProfile(ctx context.Context, accessToken string) (auth.User, error) {
	if accessToken == "" {
		return auth.User{}, provider.ErrCodeEmpty
	}

	// For reference, see:
	// https://docs.gitlab.com/ee/api/users.html#for-normal-users-1
	requestUrl := url.URL{
		Scheme: "https",
		Host:   "gitlab.com",
		Path:   "/api/v4/user",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return auth.User{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("PRIVATE-TOKEN", accessToken)

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

	userProfileUrl, err := url.Parse(responseBody.WebUrl)
	if err != nil {
		return auth.User{}, fmt.Errorf("error parsing profile url: %w", err)
	}

	userCreatedAt, err := time.Parse(time.RFC3339Nano, responseBody.CreatedAt)
	if err != nil {
		return auth.User{}, fmt.Errorf("error parsing created at: %w", err)
	}

	return auth.User{
		ID:               responseBody.ID,
		Provider:         auth.ProviderGitlab,
		NodeID:           "",
		Name:             responseBody.Name,
		Username:         responseBody.Username,
		AvatarURL:        userAvatarUrl,
		ProfileURL:       userProfileUrl,
		Location:         responseBody.Location,
		Email:            responseBody.Email,
		PublicRepository: -1,
		Followers:        responseBody.Followers,
		Following:        responseBody.Following,
		CreatedAt:        userCreatedAt,
		RegisteredAt:     time.Time{},
	}, nil
}
