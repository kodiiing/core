// auth provides authentication procedure
package auth

import (
	"context"
	"net/url"
	"time"
)

type Provider uint8

const (
	ProviderGithub Provider = iota
	ProviderGitlab
)

type User struct {
	ID       int64
	Provider Provider
	NodeID   string
	// Name refers to the display name of the user
	// that is displayed on their corresponding profile page.
	Name             string
	Username         string
	AvatarURL        *url.URL
	ProfileURL       *url.URL
	Location         string
	Email            string
	PublicRepository int64
	Followers        int64
	Following        int64
	// CreatedAt referes to the time that the user
	// register to their corresponding provider.
	CreatedAt time.Time
	// RegisteredAt refers to the time that the user is register
	// to the Kodiiing platform
	RegisteredAt time.Time
}

type Authenticate interface {
	Authenticate(ctx context.Context, accessToken string) (*User, error)
}
