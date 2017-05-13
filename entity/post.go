package entity

import "time"

type Post struct {
	Id           int64 `db:"id"`
	Key          string
	CreatedAt    *time.Time `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
	UrlName      string     `db:"url_name"`
	Title        string     `db:"title"`
	Body         string     `db:"body"`
	Labels       []string
	NextPost     *Post
	PreviousPost *Post
}
