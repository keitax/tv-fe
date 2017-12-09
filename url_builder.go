package textvid

import (
	"fmt"
	"strings"
)

// NewUrlBuilder is a constructor of UrlBuilder.
func NewUrlBuilder(conf *Config) *UrlBuilder {
	return &UrlBuilder{conf}
}

// UrlBuilder is a type to build URLs to any pages.
type UrlBuilder struct {
	config *Config
}

// LinkToTopPage builds a URL to a top page.
func (ub *UrlBuilder) LinkToTopPage() string {
	return ub.config.BaseUrl
}

// LinkToPostPage builds a URL to a single list page.
func (ub *UrlBuilder) LinkToPostPage(post *Post) string {
	return fmt.Sprintf("%s%s.html", ub.config.BaseUrl, post.Key)
}

// LinkToPostListPage builds a URL to a post list page.
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
