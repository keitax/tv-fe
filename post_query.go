package tvfe

import (
	"time"
)

// PostQuery is a value object type representing query params to get posts.
type PostQuery struct {
	Start   uint64
	Results uint64
	Year    int
	Month   time.Month
	URLName string
}

// Next makes a new post query to get next posts.
func (pq *PostQuery) Next() *PostQuery {
	nextPq := *pq
	nextPq.Start = pq.Start - pq.Results
	return &nextPq
}

// Previous makes a new post query to get previous posts.
func (pq *PostQuery) Previous() *PostQuery {
	prevPq := *pq
	prevPq.Start = pq.Start + pq.Results
	return &prevPq
}
