package main

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"fmt"
	auth_jwt "kodiiing/auth/jwt"
	auth_middleware "kodiiing/auth/middleware"
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
	"github.com/urfave/cli/v2"
)

func ApiServer(ctx context.Context) error {
	env, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		env = "development"
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		env = "5001"
	}

	// TODO: Modify this to acquire from configuration file
	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		databaseUrl = "postgres://root@localhost:5432/kodiiing?sslmode=disable"
	}

	// TODO: Modify this to acquire from configuration file
	searchUrl, ok := os.LookupEnv("SEARCH_URL")
	if !ok {
		searchUrl = "http://localhost:8108"
	}

	// TODO: Modify this to acquire from configuration file
	searchApiKey, ok := os.LookupEnv("SEARCH_API_KEY")
	if !ok {
		searchApiKey = ""
	}

	// TODO: Migrate to pgx (using pgxpool), if possible
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return fmt.Errorf("Error opening database connection: %w", err)
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

	// TODO: Make default eviction time configurable from the configuration file
	memory, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute * 3))
	if err != nil {
		return fmt.Errorf("Error creating memory cache: %w", err)
	}
	defer func(memory *bigcache.BigCache) {
		err := memory.Close()
		if err != nil {
			log.Printf("Error closing memory cache: %v", err)
		}
	}(memory)

	// TODO: Move migration to a separate command using goose (https://github.com/pressly/goose)
	//schema migration (YugaByte/PGSQL)
	errMigrateSchema := hack_provider.MigrateHackSQL(ctx, db)
	if errMigrateSchema != nil {
		return fmt.Errorf("failed to migrate: %v", errMigrateSchema)
	}
	//Collection schema (Typesense)
	errCreateCollection := hack_provider.CreateCollections(ctx, search)
	if errCreateCollection != nil {
		return fmt.Errorf("failed to migrate: %v", errMigrateSchema)
	}

	accessPublicKey, accessPrivateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("failed to generate access key pair: %v", err)
	}

	refreshPublicKey, refreshPrivateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("failed to generate refresh key pair: %v", err)
	}

	//start hacks module
	hackJwt := auth_jwt.NewJwt(accessPrivateKey, accessPublicKey, refreshPrivateKey, refreshPublicKey, "kodiiing", "user", "kodiiing")
	hackAuthService := auth_service.NewAuthService(env, db, memory)
	hackAuthMiddleware := auth_middleware.NewAuthMiddleware(hackAuthService, hackJwt)
	hackProviderTypesense := hack_provider.NewHackTypesense(search)
	hackProviderSQL := hack_provider.NewHackYugabyte(db)
	hackService := hack_service.NewHackService(env, hackAuthMiddleware, *hackProviderSQL, *hackProviderTypesense)
	//end hacks module

	app := chi.NewRouter()

	app.Mount("/Hack", hack_stub.NewHackServiceServer(hackService))
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

	return nil
}

var version string

func App() *cli.App {
	return &cli.App{
		Name:        "Kodiiing Core",
		HelpName:    "",
		Usage:       "",
		UsageText:   "",
		ArgsUsage:   "",
		Version:     version,
		Description: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "configuration-file",
				Required:  false,
				Value:     "",
				Aliases:   []string{"conf-file", "config-file"},
				EnvVars:   []string{"CONFIGURATION_FILE"},
				TakesFile: true,
			},
		},
		Copyright: `   Copyright 2023  Kodiiing

		Licensed under the Apache License, Version 2.0 (the "License");
		you may not use this file except in compliance with the License.
		You may obtain a copy of the License at

			http://www.apache.org/licenses/LICENSE-2.0

		Unless required by applicable law or agreed to in writing, software
		distributed under the License is distributed on an "AS IS" BASIS,
		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
		See the License for the specific language governing permissions and
		limitations under the License.
	 `,
		DefaultCommand: "server",
		Commands: []*cli.Command{
			{
				Name:        "server",
				Description: "Main entrypoint for Kodiiing Core. Spawns a HTTP server.",
				ArgsUsage:   "",
				Category:    "",
				Action: func(c *cli.Context) error {
					return ApiServer(c.Context)
				},
				Subcommands: []*cli.Command{},
			},
			{
				Name:        "migrate",
				Description: "Database migration",
				Subcommands: []*cli.Command{
					{
						Name: "up",
						Action: func(c *cli.Context) error {
							// TODO: Create up migration using goose (https://github.com/pressly/goose)
							return nil
						},
					},
					{
						Name: "down",
						Action: func(c *cli.Context) error {
							// TODO: Create down migration using goose (https://github.com/pressly/goose)
							return nil
						},
					},
				},
			},
		},
	}
}

func main() {
	err := App().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
