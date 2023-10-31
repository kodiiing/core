package main

import (
	"context"
	"database/sql"
	"fmt"
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
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/typesense/typesense-go/typesense"
	"github.com/urfave/cli/v2"
)

func ApiServer(ctx context.Context) error {

	config, err := GetConfig("configuration-file.yml")
	if err != nil {
		return fmt.Errorf("Error getting configuration file: %w", err)
	}

	searchUrl := fmt.Sprintf("%s:%s", config.Search.Host, config.Search.Port)

	pgxConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Databases.User,
		config.Databases.Password,
		config.Databases.Host,
		config.Databases.Port,
		config.Databases.Name,
	))
	if err != nil {
		return fmt.Errorf("error parsing database configuration: %w", err)
	}

	pgxPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pgxPool.Close()

	search := typesense.NewClient(
		typesense.WithServer(searchUrl),
		typesense.WithAPIKey(config.Search.Key),
	)

	// TODO: Make default eviction time configurable from the configuration file
	memory, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute * 3))
	if err != nil {
		return fmt.Errorf("error creating memory cache: %w", err)
	}
	defer func(memory *bigcache.BigCache) {
		err := memory.Close()
		if err != nil {
			log.Printf("Error closing memory cache: %v", err)
		}
	}(memory)

	//Collection schema (Typesense)
	errCreateCollection := hack_provider.CreateCollections(ctx, search)
	if errCreateCollection != nil {
		return fmt.Errorf("failed to migrate: %v", errCreateCollection)
	}

	app := chi.NewRouter()

	app.Mount("/Hack", hack_stub.NewHackServiceServer(hack_service.NewHackService(config.Environment, pgxPool, search)))
	app.Mount("/User", user_stub.NewUserServiceServer(user_service.NewUserService(config.Environment, pgxPool)))
	app.Mount("/Auth", auth_stub.NewAuthenticationServiceServer(auth_service.NewAuthService(config.Environment, pgxPool, memory)))
	app.Mount("/CodeReview", codereview_stub.NewCodeReviewServiceServer(codereview_service.NewCodeReviewService(config.Environment, pgxPool)))

	server := &http.Server{
		Addr:         ":" + config.Port,
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
							config, err := GetConfig(c.String("configuration-file"))
							if err != nil {
								return fmt.Errorf("Error getting configuration file: %w", err)
							}
							sql, err := sql.Open("postgres", fmt.Sprintf(
								"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
								config.Databases.User,
								config.Databases.Password,
								config.Databases.Host,
								config.Databases.Port,
								config.Databases.Name,
							))
							if err != nil {
								return fmt.Errorf("error parsing database configuration: %w", err)
							}
							defer func() {
								if err := sql.Close(); err != nil {
									log.Warn().Msgf("failed to close database connection: %v", err)
								}
							}()
							migrate, err := NewMigration(sql)
							if err != nil {
								return fmt.Errorf("failed to create migration: %w", err)
							}
							if err := migrate.Up(c.Context); err != nil {
								return fmt.Errorf("failed to migrate: %w", err)
							}
							log.Info().Msg("Migration succeed")
							return nil
						},
					},
					{
						Name: "down",
						Action: func(c *cli.Context) error {
							config, err := GetConfig(c.String("configuration-file"))
							if err != nil {
								return fmt.Errorf("Error getting configuration file: %w", err)
							}
							sql, err := sql.Open("postgres", fmt.Sprintf(
								"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
								config.Databases.User,
								config.Databases.Password,
								config.Databases.Host,
								config.Databases.Port,
								config.Databases.Name,
							))
							if err != nil {
								return fmt.Errorf("error parsing database configuration: %w", err)
							}
							defer func() {
								if err := sql.Close(); err != nil {
									log.Warn().Msgf("failed to close database connection: %v", err)
								}
							}()
							migrate, err := NewMigration(sql)
							if err != nil {
								return fmt.Errorf("failed to create migration: %w", err)
							}
							if err := migrate.Up(c.Context); err != nil {
								return fmt.Errorf("failed to migrate: %w", err)
							}
							log.Info().Msg("Migration succeed")
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
		log.Fatal().Msgf("Error running application: %v", err)
	}
}
