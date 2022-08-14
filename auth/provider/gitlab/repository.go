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

type getPublicRepositoriesResponse struct {
	ID                int64  `json:"id"`
	Description       string `json:"description"`
	Name              string `json:"name"`
	NameWithNamespace string `json:"name_with_namespace"`
	Path              string `json:"path"`
	PathWithNamespace string `json:"path_with_namespace"`
	CreatedAt         string `json:"created_at"`
	WebUrl            string `json:"web_url"`
	ForksCount        int64  `json:"forks_count"`
	StarCount         int64  `json:"star_count"`
	LastActivityAt    string `json:"last_activity_at"`
	Namespace         struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Path      string `json:"path"`
		Kind      string `json:"kind"`
		FullPath  string `json:"full_path"`
		AvatarUrl string `json:"avatar_url"`
		WebUrl    string `json:"web_url"`
	} `json:"namespace"`
}

func (g *gitlab) GetPublicRepositories(ctx context.Context, username string) ([]auth.Repository, error) {
	if username == "" {
		return []auth.Repository{}, provider.ErrCodeEmpty
	}

	// For reference, see:
	// https://docs.gitlab.com/ee/api/projects.html#list-user-projects
	requestUrl := url.URL{
		Scheme: "https",
		Host:   "gitlab.com",
		Path:   "/api/v4/users/" + username + "/projects",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return []auth.Repository{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []auth.Repository{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []auth.Repository{}, fmt.Errorf("error response: %d", resp.StatusCode)
	}

	var responseBody []getPublicRepositoriesResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return []auth.Repository{}, fmt.Errorf("error decoding response: %w", err)
	}

	repositories := make([]auth.Repository, 0, len(responseBody))
	for _, repo := range responseBody {
		repoUrl, err := url.Parse(repo.WebUrl)
		if err != nil {
			return []auth.Repository{}, fmt.Errorf("error parsing repository URL: %w", err)
		}

		repoCreatedAt, err := time.Parse(time.RFC3339Nano, repo.CreatedAt)
		if err != nil {
			return []auth.Repository{}, fmt.Errorf("error parsing repository created at: %w", err)
		}

		repoLastActivityAt, err := time.Parse(time.RFC3339Nano, repo.LastActivityAt)
		if err != nil {
			return []auth.Repository{}, fmt.Errorf("error parsing repository last activity at: %w", err)
		}

		repositories = append(repositories, auth.Repository{
			ID:             repo.ID,
			Provider:       auth.ProviderGitlab,
			Name:           repo.Name,
			Description:    repo.Description,
			URL:            repoUrl,
			Fork:           false,
			ForksCount:     repo.ForksCount,
			StarsCount:     repo.StarCount,
			OwnerUsername:  repo.Namespace.Path,
			CreatedAt:      repoCreatedAt,
			LastActivityAt: repoLastActivityAt,
		})
	}

	return repositories, nil
}
