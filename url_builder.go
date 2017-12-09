package textvid

import (
	"fmt"
	"strings"
)

// NewURLBuilder is a constructor of URLBuilder.
func NewURLBuilder(conf *Config) *URLBuilder {
	return &URLBuilder{conf}
}

// URLBuilder is a type to build URLs to any pages.
type URLBuilder struct {
	config *Config
}

// LinkToTopPage builds a URL to a top page.
func (ub *URLBuilder) LinkToTopPage() string {
	return ub.config.BaseURL
}

// LinkToPostPage builds a URL to a single list page.
func (ub *URLBuilder) LinkToPostPage(post *Post) string {
	return fmt.Sprintf("%s%s.html", ub.config.BaseURL, post.Key)
}

// LinkToPostListPage builds a URL to a post list page.
func (ub *URLBuilder) LinkToPostListPage(query *PostQuery) string {
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
	return fmt.Sprintf("%sposts/%s", ub.config.BaseURL, qs)
}
