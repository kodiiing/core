package user

import "kodiiing/auth"

type User struct {
	ID       int64
	Name     string
	Provider auth.Provider
}
