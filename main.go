package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/meilisearch/meilisearch-go"
)

func main() {
	env, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		env = "development"
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		env = "5001"
	}

	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		databaseUrl = "postgres://root@localhost:5432/kodiiing?sslmode=disable"
	}

	meilisearchUrl, ok := os.LookupEnv("MEILISEARCH_URL")
	if !ok {
		meilisearchUrl = "http://localhost:7700"
	}

	meilisearchApiKey, ok := os.LookupEnv("MEILISEARCH_API_KEY")
	if !ok {
		meilisearchApiKey = ""
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}(db)

	search := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   meilisearchUrl,
		APIKey: meilisearchApiKey,
	})

	app := chi.NewRouter()

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      app,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 15,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGINT)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("error during listening server: %v", err)
		}
	}()

	<-sig

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("error during shutting down server: %v", err)
	}
}
