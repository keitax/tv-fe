package util

import (
	"fmt"
	"strings"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/dao"
	"github.com/keitax/textvid/entity"
)

func NewUrlBuilder(conf *config.Config) *UrlBuilder {
	return &UrlBuilder{conf}
}

type UrlBuilder struct {
	config *config.Config
}

func (ub *UrlBuilder) LinkToTopPage() string {
	return ub.config.BaseUrl
}

func (ub *UrlBuilder) LinkToPostPage(post *entity.Post) string {
	return fmt.Sprintf("%s%04d/%02d/%s.html", ub.config.BaseUrl, post.CreatedAt.Year(), post.CreatedAt.Month(), post.UrlName)
}

func (ub *UrlBuilder) LinkToPostListPage(query *dao.PostQuery) string {
	q := []string{}
	if query.Start != 0 {
		q = append(q, fmt.Sprintf("start=%d", query.Start))
	}
	if query.Results != 0 {
		q = append(q, fmt.Sprintf("results=%d", query.Results))
	}
	var qs string
	if len(q) <= 0 {
		qs = ""
	} else {
		qs = "?" + strings.Join(q, "&")
	}
	return fmt.Sprintf("%sposts/%s", ub.config.BaseUrl, qs)
}
