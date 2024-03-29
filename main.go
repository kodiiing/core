package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kodiiing/telemetry"
	"kodiiing/user/user_profile"
	"net/http"
	"os"
	"os/signal"
	"time"

	authjwt "kodiiing/auth/jwt"
	authmiddleware "kodiiing/auth/middleware"
	authservice "kodiiing/auth/service"
	authstub "kodiiing/auth/stub"
	codereviewservice "kodiiing/codereview/service"
	codereviewstub "kodiiing/codereview/stub"
	hackservice "kodiiing/hack/service"
	hackstub "kodiiing/hack/stub"
	taskrepository "kodiiing/task/repository"
	taskservice "kodiiing/task/service"
	taskstub "kodiiing/task/stub"
	userservice "kodiiing/user/service"
	userstub "kodiiing/user/stub"

	hackprovider "kodiiing/hack/provider"

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
		return fmt.Errorf("getting configuration file: %w", err)
	}

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

	searchUrl := fmt.Sprintf("%s:%s", config.Search.Host, config.Search.Port)
	search := typesense.NewClient(
		typesense.WithServer(searchUrl),
		typesense.WithAPIKey(config.Search.Key),
	)

	// TODO: Make default eviction time configurable from the configuration file
	memory, err := bigcache.New(context.Background(), bigcache.DefaultConfig(time.Minute*3))
	if err != nil {
		return fmt.Errorf("error creating memory cache: %w", err)
	}
	defer func(memory *bigcache.BigCache) {
		err := memory.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Closing memory cache")
		}
	}(memory)

	//Collection schema (Typesense)
	errCreateCollection := hackprovider.CreateCollections(ctx, search)
	if errCreateCollection != nil {
		return fmt.Errorf("failed to migrate: %v", errCreateCollection)
	}

	// Build repositories
	userProfileRepository, err := user_profile.NewUserProfileRepository(pgxPool)
	if err != nil {
		return fmt.Errorf("creating user profile repository: %w", err)
	}
	taskRepository := taskrepository.NewTaskRepository(&taskrepository.Dependency{
		DB: pgxPool,
	})

	// Build service
	authService := authservice.NewAuthService(config.Environment, pgxPool, memory)

	// Build lib
	// TODO: change this into the actual value
	authJwt := authjwt.NewJwt(
		[]byte("todo_private_key"),
		[]byte("todo_public_key"),
		[]byte("todo_refresh_private_key"),
		[]byte("todo_refresh_public_key"),
		"todo_issuer",
		"todo_subject",
		"todo_audience",
	)

	// Build middleware
	authMiddleware := authmiddleware.NewAuthMiddleware(authService, authJwt)

	taskService, err := taskservice.NewTaskService(&taskservice.Config{
		Pool:           pgxPool,
		Authentication: authMiddleware,
		TaskRepository: taskRepository,
	})
	if err != nil {
		return fmt.Errorf("creating task service: %w", err)
	}

	app := chi.NewRouter()

	app.Mount("/Hack", hackstub.NewHackServiceServer(hackservice.NewHackService(config.Environment, pgxPool, search)))
	app.Mount("/User", userstub.NewUserServiceServer(userservice.NewUserService(config.Environment, userProfileRepository)))
	app.Mount("/Auth", authstub.NewAuthenticationServiceServer(authService))
	app.Mount("/CodeReview", codereviewstub.NewCodeReviewServiceServer(codereviewservice.NewCodeReviewService(config.Environment, pgxPool)))
	app.Mount("/Task", taskstub.NewTaskServiceServer(taskService))

	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      app,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 15,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// telemetry
	telemetryProvider := telemetry.NewTelemetryProvider(telemetry.Config{
		ServiceName:          "kodiiing",
		GrpcExporterEndpoint: config.Otel.ReceiverOtlpGrpcEndpoint,
		HttpExporterEndpoint: config.Otel.ReceiverOtlpHttpEndpoint,
	})

	go func() {
		log.Info().Msgf("Listening on port: %s", config.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Info().Msgf("error during listening server: %v", err)
		}
	}()

	var telemetryShutDownFuncs []telemetry.ShutDownFunc
	go func() {
		log.Info().Msg("initializing telemetry")
		telemetryShutDownFuncs, err = telemetryProvider.Run(context.Background())
		if err != nil {
			log.Fatal().Err(err).Msg("initializing telemetry")
		}
	}()

	<-sig

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("error during shutting down server: %v", err)
	}

	for _, shutDownFunc := range telemetryShutDownFuncs {
		err := shutDownFunc(shutdownCtx)
		if err != nil {
			log.Error().Err(err).Msg("shutdown telemetry...")
		}
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
								return fmt.Errorf("getting configuration file: %w", err)
							}
							db, err := sql.Open("postgres", fmt.Sprintf(
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
								if err := db.Close(); err != nil {
									log.Warn().Msgf("failed to close database connection: %v", err)
								}
							}()
							migrate, err := NewMigration(db)
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
								return fmt.Errorf("getting configuration file: %w", err)
							}
							db, err := sql.Open("postgres", fmt.Sprintf(
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
								if err := db.Close(); err != nil {
									log.Warn().Msgf("failed to close database connection: %v", err)
								}
							}()
							migrate, err := NewMigration(db)
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
		log.Fatal().Err(err).Msg("Running application")
	}
}
