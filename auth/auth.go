// Package auth provides authentication procedure
package auth

import (
	"context"
	"errors"
	"net/url"
	"time"
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
	// CreatedAt refers to the time that the user
	// register to their corresponding provider.
	CreatedAt time.Time
	// RegisteredAt refers to the time that the user is register
	// to the Kodiiing platform
	RegisteredAt time.Time
}

type Repository struct {
	ID             int64
	Provider       Provider
	Name           string
	URL            *url.URL
	Description    string
	Fork           bool
	ForksCount     int64
	StarsCount     int64
	OwnerUsername  string
	CreatedAt      time.Time
	LastActivityAt time.Time
}

type Authenticate interface {
	Authenticate(ctx context.Context, accessToken string) (*User, error)
}

// ErrUserNotFound is returned when there is a query to the database
// to find a user, yet the user was not found
var ErrUserNotFound = errors.New("user not found")

// ErrParameterEmpty is returned when a parameter is empty
// for a function call
var ErrParameterEmpty = errors.New("empty parameter")
