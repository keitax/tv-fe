package dao

import "time"

type PostQuery struct {
	Start   uint64
	Results uint64
	Year    int
	Month   time.Month
	UrlName string
}
