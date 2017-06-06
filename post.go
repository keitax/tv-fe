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

type PostList []*Post

func (sp PostList) Len() int {
	return len(sp)
}

func (sp PostList) Less(i, j int) bool {
	if sp[i].Date == nil {
		return false
	}
	if sp[j].Date == nil {
		return true
	}
	return sp[i].Date.After(*sp[j].Date)
}

func (sp PostList) Swap(i, j int) {
	sp[i], sp[j] = sp[j], sp[i]
}
