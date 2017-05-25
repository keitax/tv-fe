package textvid

import (
	"time"
)

type Post struct {
	Id           int64 `db:"id"`
	Key          string
	Date         *time.Time
	CreatedAt    *time.Time `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
	UrlName      string     `db:"url_name"`
	Title        string     `db:"title"`
	Body         string     `db:"body"`
	Labels       []string
	NextPost     *Post
	PreviousPost *Post
}

type SortPost []*Post

func (sp SortPost) Len() int {
	return len(sp)
}

func (sp SortPost) Less(i, j int) bool {
	return sp[i].Date.After(*sp[j].Date)
}

func (sp SortPost) Swap(i, j int) {
	sp[i], sp[j] = sp[j], sp[i]
}
