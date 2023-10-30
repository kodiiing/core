package hack_service

import (
	"context"
	"log"
	"net/http"
	"time"

	"kodiiing/auth"
	hack_provider "kodiiing/hack/provider"
	hack_stub "kodiiing/hack/stub"
)

type HackService struct {
	environment   string
	auth          auth.Authenticate
	HackYugabyte  hack_provider.HackYugabyte
	hackTypesense hack_provider.HackTypesense
}

func NewHackService(env string, auth auth.Authenticate, HackYugabyte hack_provider.HackYugabyte, hackTypesense hack_provider.HackTypesense) *HackService {
	return &HackService{
		environment:   env,
		auth:          auth,
		HackYugabyte:  HackYugabyte,
		hackTypesense: hackTypesense,
	}
}

// Starts a new hack post.
func (d *HackService) Create(ctx context.Context, req *hack_stub.CreateRequest) (*hack_stub.CreateResponse, *hack_stub.HackServiceError) {

	user, errAuth := d.auth.Authenticate(ctx, req.Auth.AccessToken)
	if errAuth != nil {
		return nil, &hack_stub.HackServiceError{StatusCode: http.StatusNonAuthoritativeInfo, Error: errAuth}
	}

	//store request to pgsql for main persistent storage databases
	res, err := d.HackYugabyte.CreateRepo(ctx, req)
	if err != nil {
		return &hack_stub.CreateResponse{}, &hack_stub.HackServiceError{StatusCode: http.StatusInternalServerError, Error: err}
	}

	//create document from res.Id and request, then store it to typesense for search engine databases
	document := hack_provider.HackDocument{
		Id:        res.Id,
		Title:     req.Title,
		Content:   req.Text,
		Tags:      req.Tags,
		Token:     req.Auth.AccessToken,
		Author:    hack_stub.Author{Name: user.Name, ProfileUrl: user.ProfileURL.Path, PictureUrl: user.AvatarURL.Path},
		Upvotes:   0,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}
	defer func() {
		if err := d.hackTypesense.CreateDocument(ctx, &document); err != nil {
			log.Printf("message error creating document: %v", err)
		}
	}()
	return &hack_stub.CreateResponse{Id: res.Id}, nil
}

// Upvote a hack post.
func (d *HackService) Upvote(ctx context.Context, req *hack_stub.UpvoteRequest) (*hack_stub.UpvoteResponse, *hack_stub.HackServiceError) {
	_, errAuth := d.auth.Authenticate(ctx, req.Auth.AccessToken)
	if errAuth != nil {
		return nil, &hack_stub.HackServiceError{StatusCode: http.StatusNonAuthoritativeInfo, Error: errAuth}
	}

	//store request to pgsql for main persistent storage databases
	res, err := d.HackYugabyte.UpvoteRepo(ctx, req)
	if err != nil {
		return &hack_stub.UpvoteResponse{}, &hack_stub.HackServiceError{StatusCode: http.StatusInternalServerError, Error: err}
	}
	//update document in typesense for search engine databases
	defer func() {
		if err := d.hackTypesense.UpvoteDocument(ctx, req.Id, res.Score); err != nil {
			log.Printf("message error updated document: %v", err)
		}
	}()
	return &hack_stub.UpvoteResponse{Upvoted: res.Upvoted, Score: res.Score}, nil
}

// Comment to a hack post, or reply to an existing comment.
func (d *HackService) Comment(ctx context.Context, req *hack_stub.CommentRequest) (*hack_stub.CommentResponse, *hack_stub.HackServiceError) {
	user, errAuth := d.auth.Authenticate(ctx, req.Auth.AccessToken)
	if errAuth != nil {
		return &hack_stub.CommentResponse{}, &hack_stub.HackServiceError{StatusCode: http.StatusInternalServerError, Error: errAuth}
	}
	//store request to pgsql for main persistent storage databases
	res, err := d.HackYugabyte.CommentRepo(ctx, req, user.ID)
	if err != nil {
		return &hack_stub.CommentResponse{}, &hack_stub.HackServiceError{StatusCode: http.StatusInternalServerError, Error: errAuth}
	}

	//create document from res.Id and request, then store it to typesense for search engine databases
	document := hack_stub.Author{
		Name:       user.Name,
		ProfileUrl: user.ProfileURL.Path,
		PictureUrl: user.AvatarURL.Path,
	}
	defer func() {
		if err := d.hackTypesense.CommentDocument(ctx, req, res.CommentId, &document); err != nil {
			log.Printf("message error updated document: %v", err)
		}
	}()
	return &hack_stub.CommentResponse{HackId: res.HackId, CommentId: res.CommentId}, nil
}

// See all hack posts, or maybe with a filter.
func (d *HackService) List(ctx context.Context, req *hack_stub.ListRequest) (*hack_stub.ListResponse, *hack_stub.HackServiceError) {
	_, err := d.HackYugabyte.ListRepo(ctx, req)
	if err != nil {
		return &hack_stub.ListResponse{}, &hack_stub.HackServiceError{StatusCode: http.StatusInternalServerError, Error: err}
	}

	// for i := 0; i < len(res.Hacks); i++ {
	// 	//in provider package, I added notes for author_id in the name struct of Author. (line 182)
	// 	//this's it is works, the author_id used to get the author
	// 	author, err := d.HackYugabyte.GetAuthor(ctx, res.Hacks[i].Author.Name)
	// 	if err!= nil {
	// 		return nil, &hack_stub.HackServiceError{StatusCode: http.StatusInternalServerError, Error: err}
	// 	}
	// 	comments, err := d.HackYugabyte.GetComments(ctx, res.Hacks[i].Id)
	// 	if err!= nil {
	// 		return nil, &hack_stub.HackServiceError{StatusCode: http.StatusInternalServerError, Error: err}
	// 	}
	// 	res.Hacks[i].Author.Name = author.Name
	// 	res.Hacks[i].Author.ProfileUrl = author.ProfileUrl
	// 	res.Hacks[i].Author.PictureUrl = author.PictureUrl
	// 	res.Hacks[i].Comments = comments
	// }
	//return &hack_stub.ListResponse{Hacks: res.Hacks, TotalResults: res.TotalResults, CurrentPage: res.CurrentPage, TotalPage: res.TotalPage}, nil
	return &hack_stub.ListResponse{}, nil
}
