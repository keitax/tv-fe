package tvfe

// ViewSet is a factory type to make views.
type ViewSet struct {
	urlBuilder *URLBuilder
	config     *Config
}

// NewViewSet makes a new ViewSet.
func NewViewSet(ub *URLBuilder, conf *Config) *ViewSet {
	return &ViewSet{
		urlBuilder: ub,
		config:     conf,
	}
}

// PostSingleView creates a view of the single post page.
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

// PostListView creates a view of post list page.
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
