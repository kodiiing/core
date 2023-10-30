package hack_provider

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	hack_stub "kodiiing/hack/stub"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

// go:embed ../migrations/*.sql
var embedMigrations embed.FS

type HackMigration struct {
	db *sql.DB
}

func NewHackMigration(db *sql.DB) (*HackMigration, error) {
	if db == nil {
		return &HackMigration{}, errors.New("db is nil")
	}
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return &HackMigration{}, err
	}

	return &HackMigration{db: db}, nil
}

func (m *HackMigration) Up(ctx context.Context) (err error) {
	return goose.UpContext(ctx, m.db, "../migrations")
}

func (m *HackMigration) Down(ctx context.Context) error {
	return goose.DownContext(ctx, m.db, "../migrations")
}

func CreateCollections(ctx context.Context, client *typesense.Client) *hack_stub.HackServiceError {
	// authors collection
	_, err := client.Collection("authors").Retrieve()
	if err != nil {
		log.Printf("authors collection does not exist: %v\n", err.Error())
		log.Printf("creating authors collection...\n")

		authorColletions := &api.CollectionSchema{
			Name: "authors",
			Fields: []api.Field{
				{
					Name: "name",
					Type: "string",
				},
				{
					Name: "email",
					Type: "string",
				},
				{
					Name: "access_token",
					Type: "string",
				},
				{
					Name: "profile_url",
					Type: "string",
				},
				{
					Name: "picture_url",
					Type: "string",
				},
				{
					Name: "created_at",
					Type: "string",
				},
				{
					Name: "updated_at",
					Type: "string",
				},
			},
		}
		_, err = client.Collections().Create(authorColletions)
		if err != nil {
			log.Fatal().Msgf("created authors collection error: %v", err.Error())
		}
	}

	//hacks collection
	_, err = client.Collection("hacks").Retrieve()
	if err != nil {
		log.Printf("hacks collection does not exist: %v\n", err.Error())
		log.Printf("creating hacks collection...\n")

		commentColletion := &api.CollectionSchema{
			Name: "hacks",
			Fields: []api.Field{
				{
					Name: "author",
					Type: "object",
				},
				{
					Name: "title",
					Type: "string",
				},
				{
					Name: "content",
					Type: "string",
				},
				{
					Name: "upvotes",
					Type: "int64",
				},
				{
					Name: "tags",
					Type: "string[]",
				},
				{
					Name: "comments",
					Type: "object[]",
				},
				{
					Name: "created_at",
					Type: "string",
				},
				{
					Name: "updated_at",
					Type: "string",
				},
			},
		}
		_, err = client.Collections().Create(commentColletion)
		if err != nil {
			log.Fatal().Msgf("created hacks collection error: %v", err)
		}
	}

	//comments collection
	_, err = client.Collection("comments").Retrieve()
	if err != nil {
		log.Printf("comments collection does not exist: %v\n", err.Error())
		log.Printf("creating comments collection...\n")

		commentColletion := &api.CollectionSchema{
			Name: "comments",
			Fields: []api.Field{
				{
					Name: "content",
					Type: "string",
				},
				{
					Name: "author",
					Type: "object",
				},
				{
					Name: "replies",
					Type: "object[]",
				},
				{
					Name: "created_at",
					Type: "string",
				},
			},
		}
		_, err = client.Collections().Create(commentColletion)
		if err != nil {
			log.Fatal().Msgf("created comments collection error: %v", err)
		}
	}

	log.Printf("successfully created collections...")
	return nil
}
