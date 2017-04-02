package entity

import "time"

type Post struct {
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	UrlName   string
	Labels    []string
	Title     string
	Body      string
}
