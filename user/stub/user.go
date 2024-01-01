package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserServiceError struct {
	StatusCode int
	Error      error
}

type OnboardingRequest struct {
	Reason      JoinReason     `json:"reason"`
	ReasonOther string         `json:"reason_other"`
	CodedBefore bool           `json:"coded_before"`
	Languages   []string       `json:"languages"`
	Target      string         `json:"target"`
	Auth        Authentication `json:"auth"`
}

type EmptyResponse struct {
}

type Authentication struct {
	AccessToken string `json:"access_token"`
}

type JoinReason uint32

const (
	JoinReasonJOIN_REASON_UNSPECIFIED JoinReason = 0
	JoinReasonJOIN_REASON_SCHOOL      JoinReason = 1
	JoinReasonJOIN_REASON_WORK        JoinReason = 2
	JoinReasonJOIN_REASON_FASCINATED  JoinReason = 3
	JoinReasonJOIN_REASON_FRIEND      JoinReason = 4
	JoinReasonJOIN_REASON_OTHER       JoinReason = 5
)

type UserServiceServer interface {
	Onboarding(ctx context.Context, req *OnboardingRequest) (*EmptyResponse, *UserServiceError)
}

func NewUserServiceServer(implementation UserServiceServer) *chi.Mux {
	mux := chi.NewMux()
	mux.Post("/Onboarding", func(w http.ResponseWriter, r *http.Request) {
		var req OnboardingRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[UserService - Onboardingerror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.Onboarding(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[UserService - Onboardingerror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[UserService - Onboardingerror] writing to response stream: %s", e.Error())
		}
	})

	return mux
}
