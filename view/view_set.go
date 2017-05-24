package view

import (
	"github.com/keitax/textvid"
	"github.com/keitax/textvid/entity"
	"github.com/keitax/textvid/repository"
	"github.com/keitax/textvid/urlbuilder"
)

type ViewSet struct {
	urlBuilder *urlbuilder.UrlBuilder
	config     *textvid.Config
}

func NewViewSet(ub *urlbuilder.UrlBuilder, conf *textvid.Config) *ViewSet {
	return &ViewSet{
		urlBuilder: ub,
		config:     conf,
	}
}

func (vs *ViewSet) PostSingleView(p *entity.Post) *View {
	return &View{
		vs.urlBuilder,
		vs.config,
		"post_single.tmpl",
		map[string]interface{}{
			"post": p,
		},
	}
}

func (vs *ViewSet) PostListView(ps []*entity.Post, nextPosts []*entity.Post, previousPosts []*entity.Post, q *repository.PostQuery) *View {
	return &View{
		vs.urlBuilder,
		vs.config,
		"post_list.tmpl",
		map[string]interface{}{
			"posts":         ps,
			"NextPosts":     nextPosts,
			"PreviousPosts": previousPosts,
			"CurrentQuery":  q,
		},
	}
}

func (vs *ViewSet) AdminView(ps []*entity.Post) *View {
	return &View{
		vs.urlBuilder,
		vs.config,
		"admin.tmpl",
		map[string]interface{}{
			"posts": ps,
		},
	}
}

func (vs *ViewSet) PostEditorView(p *entity.Post) *View {
	return &View{
		vs.urlBuilder,
		vs.config,
		"post_editor.tmpl",
		map[string]interface{}{
			"post": p,
		},
	}
}
