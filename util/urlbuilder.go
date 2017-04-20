package util

import (
	"fmt"

	"github.com/keitax/textvid/config"
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
