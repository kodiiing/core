package hack_provider

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	hack_stub "kodiiing/hack/stub"
	"log"

	"github.com/pressly/goose/v3"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

//go:embed migrations/*.sql
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
		return &HackMigration{}, fmt.Errorf("failed to set dialect: %s", err.Error())
	}

	return &HackMigration{db: db}, nil
}

func (m *HackMigration) Up(ctx context.Context) (err error) {
	return goose.UpContext(ctx, m.db, ".")
}

func (m *HackMigration) Down(ctx context.Context) (err error) {
	return goose.DownContext(ctx, m.db, ".")
}

func MigrateHackSQL(ctx context.Context, d *sql.DB) *hack_stub.HackServiceError {
	db, err := d.Conn(ctx)
	if err != nil {
		return &hack_stub.HackServiceError{StatusCode: 500, Error: fmt.Errorf("message err failed to connect to database:  %s", err.Error())}
	}
	defer func() {
		err := db.Close()
		if err != nil && !errors.Is(err, sql.ErrConnDone) {
			log.Printf("failed to close database connection: %V", err)
		}
	}()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return &hack_stub.HackServiceError{StatusCode: 500, Error: fmt.Errorf("failed to begin transaction: %s", err.Error())}
	}
	defer func() {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("rollback error : %s", errRollback.Error())
		}
	}()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS authors(
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			access_token TEXT CONSTRAINT unique_access_token UNIQUE NOT NULL,
			profile_url VARCHAR(255) NOT NULL,
			picture_url VARCHAR(255) NOT NULL,
			created_at timestamp default current_timestamp,
			updated_at timestamp default current_timestamp
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id BIGSERIAL PRIMARY KEY,
			content TEXT NOT NULL,
			author_id BIGSERIAL NOT NULL,
			created_at timestamp default current_timestamp,
			CONSTRAINT fk_author_comments foreign key (author_id) REFERENCES authors(id)
		)`,
		`CREATE TABLE IF NOT EXISTS hacks(
			id BIGSERIAL PRIMARY KEY,
			author_id BIGSERIAL NOT NULL,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			tags VARCHAR[] NOT NULL,
			upvotes BIGINT null,
			created_at timestamp default current_timestamp,
			updated_at timestamp default current_timestamp,
			CONSTRAINT fk_author_comments foreign key (author_id) REFERENCES authors(id)
		)`,
		`CREATE TABLE IF NOT EXISTS hack_comments (
			id BIGSERIAL PRIMARY KEY,
			hack_id BIGSERIAL NOT NULL,
			comment_id BIGSERIAL NOT NULL,
			parent_id BIGINT null,
			created_at timestamp default current_timestamp,
			CONSTRAINT fk_hack foreign key (hack_id) REFERENCES hacks(id),
			CONSTRAINT fk_comment foreign key (comment_id) REFERENCES comments(id)
		)`,
	}
	for _, query := range queries {
		_, err := tx.ExecContext(ctx, query)
		if err != nil {
			return &hack_stub.HackServiceError{StatusCode: 500, Error: fmt.Errorf("migration error: %v", err.Error())}
		}
	}

	if errCommit := tx.Commit(); errCommit != nil {
		return &hack_stub.HackServiceError{StatusCode: 500, Error: err}
	}

	log.Printf("Migration successfully!")
	return nil

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
			log.Fatalf("created authors collection error: %v", err.Error())
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
			log.Fatalf("created hacks collection error: %v", err)
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
			log.Fatalf("created comments collection error: %v", err)
		}
	}

	log.Printf("successfully created collections...")
	return nil
}
