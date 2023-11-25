package user_profile

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	user_stub "kodiiing/user/stub"
	"time"
)

type UserProfile struct {
	UserID          int64
	JoinReason      user_stub.JoinReason
	JoinReasonOther string
	HasCodedBefore  bool
	Languages       []string
	Target          string
	UpdatedAt       time.Time
	UpdatedBy       string
}

type Repository struct {
	db *pgxpool.Pool
}

func NewUserProfileRepository(db *pgxpool.Pool) (*Repository, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	return &Repository{db: db}, nil
}
