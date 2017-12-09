package textvid

import (
	"time"
)

// PostQuery is a value object type representing query params to get posts.
type PostQuery struct {
	Start   uint64
	Results uint64
	Year    int
	Month   time.Month
	UrlName string
}

// Next makes a new post query to get next posts.
func (pq *PostQuery) Next() *PostQuery {
	pq_ := *pq
	pq_.Start = pq.Start - pq.Results
	return &pq_
}

// Previous makes a new post query to get previous posts.
func (pq *PostQuery) Previous() *PostQuery {
	pq_ := *pq
	pq_.Start = pq.Start + pq.Results
	return &pq_
}
