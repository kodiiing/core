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

type getPublicRepositoriesResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login     string `json:"login"`
		ID        int64  `json:"id"`
		AvatarUrl string `json:"avatar_url"`
		HtmlUrl   string `json:"html_url"`
		Type      string `json:"type"`
	} `json:"owner"`
	HtmlUrl         string `json:"html_url"`
	Description     string `json:"description"`
	Fork            bool   `json:"fork"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	StargazersCount int64  `json:"stargazers_count"`
	Language        string `json:"language"`
	Forks           int64  `json:"forks"`
	DefaultBranch   string `json:"default_branch"`
}

func (g *github) GetPublicRepositories(ctx context.Context, username string) ([]auth.Repository, error) {
	if username == "" {
		return []auth.Repository{}, provider.ErrCodeEmpty
	}

	// For reference, see:
	// https://docs.github.com/en/rest/repos/repos#list-repositories-for-a-user
	requestUrl := url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "/users/" + username + "/repos",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return []auth.Repository{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []auth.Repository{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []auth.Repository{}, fmt.Errorf("error getting public repositories: %s", resp.Status)
	}

	var responseBody []getPublicRepositoriesResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return []auth.Repository{}, fmt.Errorf("error decoding response: %w", err)
	}

	repositories := make([]auth.Repository, 0, len(responseBody))
	for _, repo := range responseBody {
		repoUrl, err := url.Parse(repo.HtmlUrl)
		if err != nil {
			return []auth.Repository{}, fmt.Errorf("error parsing repository url: %w", err)
		}

		repoCreatedAt, err := time.Parse(time.RFC3339, repo.CreatedAt)
		if err != nil {
			return []auth.Repository{}, fmt.Errorf("error parsing repository created at: %w", err)
		}

		repoLastActivityAt, err := time.Parse(time.RFC3339, repo.UpdatedAt)
		if err != nil {
			return []auth.Repository{}, fmt.Errorf("error parsing repository last activity at: %w", err)
		}

		repositories = append(repositories, auth.Repository{
			ID:             repo.ID,
			Provider:       auth.ProviderGithub,
			Name:           repo.Name,
			URL:            repoUrl,
			Description:    repo.Description,
			Fork:           repo.Fork,
			ForksCount:     repo.Forks,
			StarsCount:     repo.StargazersCount,
			OwnerUsername:  repo.Owner.Login,
			CreatedAt:      repoCreatedAt,
			LastActivityAt: repoLastActivityAt,
		})
	}

	return repositories, nil
}
