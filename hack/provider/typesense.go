package hack_provider

import (
	"context"
	"fmt"
	hack_stub "kodiiing/hack/stub"
	"time"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

// HackDocument document represents a request to store on Typesense.
type HackDocument struct {
	Id        string           `json:"id"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	Tags      []string         `json:"tags"`
	Token     string           `json:"token"`
	Author    hack_stub.Author `json:"author"`
	Upvotes   int64            `json:"upvotes"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
}

type HackTypesense struct {
	search *typesense.Client
}

func NewHackTypesense(search *typesense.Client) *HackTypesense {
	return &HackTypesense{search: search}
}

func (d *HackTypesense) CreateDocument(ctx context.Context, req *HackDocument) error {
	document := HackDocument{
		Id:        req.Id,
		Title:     req.Title,
		Content:   req.Content,
		Tags:      req.Tags,
		Token:     req.Token,
		Author:    req.Author,
		Upvotes:   req.Upvotes,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}

	_, err := d.search.Collection("hacks").Documents().Create(document)
	if err != nil {
		return err
	}
	return nil
}

func (d *HackTypesense) UpvoteDocument(ctx context.Context, id string, score int64) error {
	document := HackDocument{
		Upvotes:   score,
		UpdatedAt: time.Now().String(),
	}
	_, err := d.search.Collection("hacks").Document(id).Update(document)
	if err != nil {
		return err
	}
	return nil
}

func (d *HackTypesense) CommentDocument(ctx context.Context, req *hack_stub.CommentRequest, commentId string, author *hack_stub.Author) error {

	_, err := d.search.Collection("hacks").Document(req.HackId).Retrieve()
	if err != nil {
		return err
	}

	document := hack_stub.Comment{
		Id:      commentId,
		Content: req.Text,
		Author: hack_stub.Author{
			Name:       author.Name,
			ProfileUrl: author.ProfileUrl,
			PictureUrl: author.PictureUrl,
		},
		CreatedAt: time.Now().String(),
	}
	_, err = d.search.Collection("comments").Documents().Create(document)
	if err != nil {
		return err
	}
	return nil
}

func (d *HackTypesense) ListDocuments(ctx context.Context, req *hack_stub.ListRequest) (*hack_stub.ListResponse, error) {
	page := int(req.Page)
	sortBy := fmt.Sprint(req.SortBy)
	perPage := 100

	//d.search.Collection("sdsd")
	data, err := d.search.Collection("hacks").Documents().Search(&api.SearchCollectionParams{Page: &page, SortBy: &sortBy, PerPage: &perPage})
	if err != nil {
		return nil, err
	}

	// hacks, err := d.search.Collection("hacks").Retrieve()

	return &hack_stub.ListResponse{CurrentPage: uint32(*data.Page), TotalPage: uint32(*data.OutOf), TotalResults: uint64(*data.Found)}, nil
}
