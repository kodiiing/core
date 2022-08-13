package auth_stub

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthenticationServiceError struct {
	StatusCode int
	Error      error
}

type LoginRequest struct {
	Provider   Provider `json:"provider"`
	AccessCode string   `json:"access_code"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type LogoutRequest struct {
	AccessToken string `json:"access_token"`
}

type EmptyResponse struct {
}

type Provider uint32

const (
	ProviderGITHUB Provider = 0
	ProviderGITLAB Provider = 1
)

type AuthenticationServiceServer interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, *AuthenticationServiceError)
	Logout(ctx context.Context, req *LogoutRequest) (*EmptyResponse, *AuthenticationServiceError)
}

func NewAuthenticationServiceServer(implementation AuthenticationServiceServer) *chi.Mux {
	mux := chi.NewMux()
	mux.Post("/Login", func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[AuthenticationService - Loginerror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.Login(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[AuthenticationService - Loginerror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[AuthenticationService - Loginerror] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/Logout", func(w http.ResponseWriter, r *http.Request) {
		var req LogoutRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[AuthenticationService - Logouterror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.Logout(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[AuthenticationService - Logouterror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[AuthenticationService - Logouterror] writing to response stream: %s", e.Error())
		}
	})

	return mux
}
