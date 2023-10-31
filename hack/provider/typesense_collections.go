package hack_provider

import (
	"context"
	hack_stub "kodiiing/hack/stub"

	"github.com/rs/zerolog/log"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

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
