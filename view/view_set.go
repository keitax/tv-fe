package view

import (
	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/entity"
	"github.com/keitax/textvid/util"
)

type ViewSet struct {
	urlBuilder *util.UrlBuilder
	config     *config.Config
}

func NewViewSet(ub *util.UrlBuilder, conf *config.Config) *ViewSet {
	return &ViewSet{
		urlBuilder: ub,
		config:     conf,
	}
}

func (vs *ViewSet) PostSingleView(p *entity.Post) View_ {
	return &view{
		vs.urlBuilder,
		vs.config,
		"post_single.tmpl",
		map[string]interface{}{
			"post": p,
		},
	}
}
