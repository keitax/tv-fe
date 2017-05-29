package textvid

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PostController struct {
	repository *Repository
	viewSet    *ViewSet
	urlBuilder *UrlBuilder
	config     *Config
}

func NewPostController(r *Repository, vs *ViewSet, ub *UrlBuilder, config_ *Config) *PostController {
	return &PostController{
		r,
		vs,
		ub,
		config_,
	}
}

func (pc *PostController) GetIndex(w http.ResponseWriter, req *http.Request) {
	q := &PostQuery{
		Start:   1,
		Results: 5,
	}
	posts := pc.repository.Fetch(q)
	qp := q.Previous()
	qp.Results = 1
	prevPosts := pc.repository.Fetch(qp)
	pc.viewSet.PostListView(posts, nil, prevPosts, q).Render(w)
}

func (pc *PostController) GetSingle(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	key := fmt.Sprintf("%s/%s/%s", params["year"], params["month"], params["name"])
	p := pc.repository.FetchOne(key)
	pc.viewSet.PostSingleView(p).Render(w)
}

func (pc *PostController) GetList(w http.ResponseWriter, req *http.Request) {
	s, err := strconv.Atoi(req.URL.Query().Get("start"))
	if err != nil {
		panic(err)
	}
	r, err := strconv.Atoi(req.URL.Query().Get("results"))
	if err != nil {
		panic(err)
	}
	q := &PostQuery{
		Start:   uint64(s),
		Results: uint64(r),
	}
	ps := pc.repository.Fetch(q)
	nextPosts := pc.repository.Fetch(q.Next())
	prevPosts := pc.repository.Fetch(q.Previous())
	pc.viewSet.PostListView(ps, nextPosts, prevPosts, q).Render(w)
}

func (pc *PostController) GetEditor(w http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	var p *Post
	if key == "" {
		p = &Post{}
	} else {
		p = pc.repository.FetchOne(key)
		if p == nil {
			http.NotFound(w, req)
			return
		}
	}
	pc.viewSet.PostEditorView(p).Render(w)
}

func (pc *PostController) SubmitPost(w http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	if err := req.ParseForm(); err != nil {
		panic(err)
	}
	pc.repository.Commit(&Post{
		Key:     key,
		Title:   req.Form.Get("title"),
		Body:    req.Form.Get("body"),
		UrlName: req.Form.Get("url-name"),
	})
	committed := pc.repository.FetchOne(key)
	http.Redirect(w, req, pc.urlBuilder.LinkToPostPage(committed), http.StatusSeeOther)
}

func (pc *PostController) EditPost(w http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	if err := req.ParseForm(); err != nil {
		panic(err)
	}
	pc.repository.Commit(&Post{
		Key:     key,
		Title:   req.Form.Get("title"),
		Body:    req.Form.Get("body"),
		UrlName: req.Form.Get("url-name"),
	})
	committed := pc.repository.FetchOne(req.Form.Get("key"))
	http.Redirect(w, req, pc.urlBuilder.LinkToPostPage(committed), http.StatusSeeOther)
}