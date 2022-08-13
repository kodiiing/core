package hack

import (
	"net/url"
	"time"
)

type Hack struct {
	Id        string
	Title     string
	Content   string
	Upvotes   int64
	Author    Author
	Comments  []Comment
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Author struct {
	Name       string
	ProfileUrl *url.URL
	PictureUrl *url.URL
}

type Comment struct {
	Id        string
	Content   string
	Author    Author
	Replies   []Comment
	CreatedAt time.Time
}
