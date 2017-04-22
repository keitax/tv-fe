package dao

import (
	"time"
)

type PostQuery struct {
	Start   uint64
	Results uint64
	Year    int
	Month   time.Month
	UrlName string
}

func (pq *PostQuery) Next() *PostQuery {
	pq_ := *pq
	pq_.Start = pq.Start - pq.Results
	return &pq_
}
