package hack_provider

import (
	"context"
	"fmt"
	hackstub "kodiiing/hack/stub"
	"log"
	"math"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type HackYugabyte struct {
	pool *pgxpool.Pool
}

func NewHackYugabyte(pool *pgxpool.Pool) *HackYugabyte {
	return &HackYugabyte{pool: pool}
}

// CreateRepo Starts a new hack post.
func (d *HackYugabyte) CreateRepo(ctx context.Context, req *hackstub.CreateRequest) (*hackstub.CreateResponse, error) {
	db, err := d.pool.Acquire(ctx)
	if err != nil {
		return &hackstub.CreateResponse{}, fmt.Errorf("message err failed to connect to database:  %s", err.Error())
	}
	defer db.Release()

	tx, err := db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return &hackstub.CreateResponse{}, fmt.Errorf("message err failed to start transaction: %s", err.Error())
	}

	defer func() {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			log.Printf("message err rollback create new hack post: %v", err.Error())
		}
	}()

	var lastInsertId string
	query := "insert into hacks (title, content, tags, access_token) values($1, $2, $3, $4) RETURNING id"
	errQueryExceHack := tx.QueryRow(ctx, query, req.Title, req.Text, req.Tags, req.Auth.AccessToken).Scan(&lastInsertId)
	if errQueryExceHack != nil {
		return &hackstub.CreateResponse{}, fmt.Errorf("message err query failed to create a new hack: %v", errQueryExceHack)
	}

	if err := tx.Commit(ctx); err != nil {
		return &hackstub.CreateResponse{}, fmt.Errorf("message err committing transaction: %w", err)
	}
	return &hackstub.CreateResponse{Id: lastInsertId}, nil
}

// UpvoteRepo Upvote a hack post.
func (d *HackYugabyte) UpvoteRepo(ctx context.Context, req *hackstub.UpvoteRequest) (*hackstub.UpvoteResponse, error) {
	db, err := d.pool.Acquire(ctx)
	if err != nil {
		return &hackstub.UpvoteResponse{}, fmt.Errorf("message err failed to connect to database:  %s", err.Error())
	}
	defer db.Release()

	tx, err := db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return &hackstub.UpvoteResponse{}, fmt.Errorf("message err failed to start transaction: %s", err.Error())
	}
	defer func() {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			log.Printf("message err rollback create new hack post: %v", err.Error())
		}
	}()

	// var isExists bool
	// query := "select exists(select hacks.id from hacks where hacks.id = $1)"
	// err = tx.QueryRowContext(ctx, query, req.Id).Scan(&isExists)
	// if err != nil {
	// 	return &hackstub.UpvoteResponse{}, &hackstub.HackServiceError{StatusCode: 500, Error: fmt.Errorf("message err failed to scan data not found: %v", err)}
	// }

	query := "select upvotes from hacks where id = $1"
	var upvote int64
	err = tx.QueryRow(ctx, query, req.Id).Scan(&upvote)
	if err != nil {
		return &hackstub.UpvoteResponse{}, fmt.Errorf("message err failed to scan data not found: %v", err)
	}
	query = "UPDATE hacks SET upvotes = $1 where id = $2"
	_, err = tx.Exec(ctx, query, upvote+1, req.Id)
	if err != nil {
		return &hackstub.UpvoteResponse{}, fmt.Errorf("message err failed to execute update query: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return &hackstub.UpvoteResponse{}, fmt.Errorf("message err committing transaction: %w", err)
	}

	return &hackstub.UpvoteResponse{Upvoted: true, Score: upvote + 1}, nil
}

// CommentRepo Comment to a hack post, or reply to an existing comment.
func (d *HackYugabyte) CommentRepo(ctx context.Context, req *hackstub.CommentRequest, authorId int64) (*hackstub.CommentResponse, error) {
	db, err := d.pool.Acquire(ctx)
	if err != nil {
		return &hackstub.CommentResponse{}, fmt.Errorf("message err failed to connect to database:  %s", err.Error())
	}
	defer db.Release()

	tx, err := db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return &hackstub.CommentResponse{}, fmt.Errorf("message err failed to start transaction: %s", err.Error())
	}
	defer func() {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			log.Printf("message err rollback create new hack post: %v", err.Error())
		}
	}()

	var commentId string
	var lastInsertId string
	query := "insert into comments (content, author_id) values ($1, $2) RETURNING id"
	errQueryComment := tx.QueryRow(ctx, query, req.Text, authorId).Scan(&lastInsertId)
	if errQueryComment != nil {
		return &hackstub.CommentResponse{}, fmt.Errorf("message err failed to execute insert comment query: %w", err)
	}

	//check if comment is commentary on a comment reply
	commentId = lastInsertId
	if req.ParentId != "" {
		lastInsertId = req.ParentId
	}

	query = "insert into hack_comments (hack_id, comment_id) values ($1, $2)"
	errQueryHackComment := tx.QueryRow(ctx, query, req.HackId, lastInsertId)
	if errQueryHackComment != nil {
		return &hackstub.CommentResponse{}, fmt.Errorf("message err failed to execute insert hack_comments query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return &hackstub.CommentResponse{}, fmt.Errorf("message err committing transaction: %w", err)
	}
	return &hackstub.CommentResponse{HackId: req.HackId, CommentId: commentId}, nil
}

// See all hack posts, or maybe with a filter.
func (d *HackYugabyte) ListRepo(ctx context.Context, req *hackstub.ListRequest) (*hackstub.ListResponse, error) {
	db, err := d.pool.Acquire(ctx)
	if err != nil {
		return &hackstub.ListResponse{}, fmt.Errorf("message err failed to connect to database:  %s", err.Error())
	}
	defer db.Release()

	tx, err := db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return &hackstub.ListResponse{}, fmt.Errorf("message err failed to start transaction: %s", err.Error())
	}
	defer func() {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			log.Printf("message err rollback list hack posts: %v", err.Error())
		}
	}()

	page := req.Page
	if page == 0 {
		page = 1
	}

	//note, author_id store to name just for a temporary. in the service, it will be replaced with the author's name.
	//why does this shit happen? 'cuz I can't add the new struct for the repo response :D.
	// query := "select count(*) over() as total, author_id, id, title, content, upvotes from hacks ORDER BY $1 ASC LIMIT 10 OFFSET 10 * ($2 -1)"

	//get hacks
	query := "select count(*) over() as total_result, id as hack_id, id, title, content, upvotes, author_id from hacks ORDER BY $1 ASC LIMIT 10 OFFSET 10 * ($2 -1)"
	rows, err := tx.Query(ctx, query, req.SortBy, page)
	if err != nil {
		return &hackstub.ListResponse{}, fmt.Errorf("message err failed to query: %v", err)
	}
	defer rows.Close()

	var hacks []hackstub.Hack
	var totalResult uint64
	var authorIds []string
	var hackIds []string
	mapHackIdx := make(map[int]int)
	idx := 0
	for rows.Next() {
		var hack hackstub.Hack
		var authorId string
		var hackId string

		if errScan := rows.Scan(&totalResult, &hackId, &hack.Id, &hack.Title, &hack.Content, &hack.Upvotes, &authorId); errScan != nil {
			return &hackstub.ListResponse{}, fmt.Errorf("message err failed to scan hack row: %v", errScan)
		}
		hacks = append(hacks, hack)
		authorIds = append(authorIds, authorId)
		hackIds = append(hackIds, hackId)
		iterate, err := strconv.Atoi(hack.Id)
		if err != nil {
			return &hackstub.ListResponse{}, fmt.Errorf("message err failed to convert hackId to int: %v", err)
		}
		mapHackIdx[iterate] = idx
		idx++
	}

	//get author
	query = "select id,name, profile_url, picture_url from authors where id = any($1)"
	rowAuthor, err := tx.Query(ctx, query, pq.Array(authorIds))
	if err != nil {
		return &hackstub.ListResponse{}, fmt.Errorf("message err failed to scan author row: %v", err)
	}
	defer rowAuthor.Close()
	for rowAuthor.Next() {
		var id int
		var author hackstub.Author
		if errScan := rowAuthor.Scan(&id, &author.Name, &author.ProfileUrl, &author.PictureUrl); errScan != nil {
			return &hackstub.ListResponse{}, fmt.Errorf("message err failed to scan hack row: %v", errScan)
		}
		i := mapHackIdx[id]
		hacks[i].Author.Name = author.Name
		hacks[i].Author.ProfileUrl = author.ProfileUrl
		hacks[i].Author.PictureUrl = author.PictureUrl
	}

	query = `select hacks.id, c.id, c.content, a.name, a.profile_url, a.picture_url, c.created_at from hacks
				left join hack_comments hc on hacks.id = hc.hack_id
				left join comments c on hc.comment_id = c.id
				left join authors a on c.author_id = a.id
				where hacks.id = any($1) order by hacks.id;`
	rowComments, err := tx.Query(ctx, query, pq.Array(hackIds))
	if err != nil {
		return &hackstub.ListResponse{}, fmt.Errorf("message err failed to query: %v", err)
	}

	for rowComments.Next() {
		var comment hackstub.Comment
		var hack_id int
		err = rowComments.Scan(&hack_id, &comment.Id, &comment.Content, &comment.Author.Name, &comment.Author.ProfileUrl, &comment.Author.PictureUrl, &comment.CreatedAt)
		if err != nil {
			return &hackstub.ListResponse{}, fmt.Errorf("message err failed to scan comment row: %v", err)
		}

		i := mapHackIdx[hack_id]

		if hacks[i].Comments == nil {
			hacks[i].Comments = make([]hackstub.Comment, 0)
		}
		hacks[i].Comments = append(hacks[i].Comments, comment)
	}
	//total page
	totalPage := math.Round(math.Ceil(float64(totalResult) / 10))

	if err := tx.Commit(ctx); err != nil {
		return &hackstub.ListResponse{}, fmt.Errorf("message err committing transaction: %w", err)
	}
	return &hackstub.ListResponse{Hacks: hacks, TotalResults: totalResult, CurrentPage: req.Page, TotalPage: uint32(totalPage)}, nil
}
