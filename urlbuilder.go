package textvid

import (
	"fmt"
	"strings"

	"github.com/keitax/textvid/entity"
	"github.com/keitax/textvid/repository"
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

func (ub *UrlBuilder) LinkToPostPage(post *entity.Post) string {
	return fmt.Sprintf("%s%s.html", ub.config.BaseUrl, post.Key)
}

func (ub *UrlBuilder) LinkToPostResource(post *entity.Post) string {
	var path string
	if post == nil {
		path = "posts/"
	} else {
		path = fmt.Sprintf("posts/%d", post.Id)
	}
	return ub.config.BaseUrl + path
}

func (ub *UrlBuilder) LinkToPostListPage(query *repository.PostQuery) string {
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

func (ub *UrlBuilder) LinkToPostEditorPage(post *entity.Post) string {
	if post == nil {
		return fmt.Sprintf("%sposts/new", ub.config.BaseUrl)
	}
	return fmt.Sprintf("%sposts/%d/edit", ub.config.BaseUrl, post.Id)
}
