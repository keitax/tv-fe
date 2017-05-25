package textvid

type ViewSet struct {
	urlBuilder *UrlBuilder
	config     *Config
}

func NewViewSet(ub *UrlBuilder, conf *Config) *ViewSet {
	return &ViewSet{
		urlBuilder: ub,
		config:     conf,
	}
}

func (vs *ViewSet) PostSingleView(p *Post) *View {
	return &View{
		vs.urlBuilder,
		vs.config,
		"post_single.tmpl",
		map[string]interface{}{
			"post": p,
		},
	}
}

func (vs *ViewSet) PostListView(ps []*Post, nextPosts []*Post, previousPosts []*Post, q *PostQuery) *View {
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

func (vs *ViewSet) AdminView(ps []*Post) *View {
	return &View{
		vs.urlBuilder,
		vs.config,
		"admin.tmpl",
		map[string]interface{}{
			"posts": ps,
		},
	}
}

func (vs *ViewSet) PostEditorView(p *Post) *View {
	return &View{
		vs.urlBuilder,
		vs.config,
		"post_editor.tmpl",
		map[string]interface{}{
			"post": p,
		},
	}
}
