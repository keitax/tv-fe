package textvid

import (
	"fmt"
	"strings"
)

func NewUrlBuilder(conf *Config) *UrlBuilder {
	return &UrlBuilder{conf}
}

type UrlBuilder struct {
	config *Config
}

func (ub *UrlBuilder) LinkToTopPage() string {
	return ub.config.BaseUrl
}

func (ub *UrlBuilder) LinkToPostPage(post *Post) string {
	return fmt.Sprintf("%s%s.html", ub.config.BaseUrl, post.Key)
}

func (ub *UrlBuilder) LinkToPostListPage(query *PostQuery) string {
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
