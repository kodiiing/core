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

	auth_service "kodiiing/auth/service"
	auth_stub "kodiiing/auth/stub"
	codereview_service "kodiiing/codereview/service"
	codereview_stub "kodiiing/codereview/stub"
	hack_provider "kodiiing/hack/provider"
	hack_service "kodiiing/hack/service"
	hack_stub "kodiiing/hack/stub"
	user_service "kodiiing/user/service"
	user_stub "kodiiing/user/stub"

	"github.com/allegro/bigcache/v3"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/typesense/typesense-go/typesense"
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

	searchUrl, ok := os.LookupEnv("SEARCH_URL")
	if !ok {
		searchUrl = "http://localhost:8108"
	}

	searchApiKey, ok := os.LookupEnv("SEARCH_API_KEY")
	if !ok {
		searchApiKey = ""
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}(db)

	search := typesense.NewClient(
		typesense.WithServer(searchUrl),
		typesense.WithAPIKey(searchApiKey),
	)

	memory, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute * 3))
	if err != nil {
		log.Fatalf("Error creating memory cache: %v", err)
		return
	}
	defer func(memory *bigcache.BigCache) {
		err := memory.Close()
		if err != nil {
			log.Printf("Error closing memory cache: %v", err)
		}
	}(memory)

	//context
	ctx := context.TODO()
	//schema migration (YugaByte/PGSQL)
	errMigrateSchema := hack_provider.MigrateHackSQL(ctx, db)
	if errMigrateSchema != nil {
		log.Fatalf("failed to migrate: %v", errMigrateSchema)
	}
	//Collection schema (Typesense)
	errCreateCollection := hack_provider.CreateCollections(ctx, search)
	if errCreateCollection != nil {
		log.Fatalf("failed to migrate: %v", errMigrateSchema)
	}

	app := chi.NewRouter()

	app.Mount("/Hack", hack_stub.NewHackServiceServer(hack_service.NewHackService(env, db, search)))
	app.Mount("/User", user_stub.NewUserServiceServer(user_service.NewUserService(env, db)))
	app.Mount("/Auth", auth_stub.NewAuthenticationServiceServer(auth_service.NewAuthService(env, db, memory)))
	app.Mount("/CodeReview", codereview_stub.NewCodeReviewServiceServer(codereview_service.NewCodeReviewService(env, db)))

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
