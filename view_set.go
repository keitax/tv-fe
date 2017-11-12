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
			"stylesheets": []string{"index.css"},
			"post":        p,
		},
	}
}

func (vs *ViewSet) PostListView(ps []*Post, nextPosts []*Post, previousPosts []*Post, q *PostQuery) *View {
	return &View{
		vs.urlBuilder,
		vs.config,
		"post_list.tmpl",
		map[string]interface{}{
			"stylesheets":   []string{"index.css"},
			"posts":         ps,
			"NextPosts":     nextPosts,
			"PreviousPosts": previousPosts,
			"CurrentQuery":  q,
		},
	}
}
