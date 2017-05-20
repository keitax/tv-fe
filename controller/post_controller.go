package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/entity"
	"github.com/keitax/textvid/repository"
	"github.com/keitax/textvid/urlbuilder"
	"github.com/keitax/textvid/view"
)

type PostController struct {
	repository *repository.Repository
	viewSet    *view.ViewSet
	urlBuilder *urlbuilder.UrlBuilder
	config     *config.Config
}

func NewPostController(r *repository.Repository, vs *view.ViewSet, ub *urlbuilder.UrlBuilder, config_ *config.Config) *PostController {
	return &PostController{
		r,
		vs,
		ub,
		config_,
	}
}

func (c *PostController) GetIndex(w http.ResponseWriter, req *http.Request) {
	q := &repository.PostQuery{
		Start:   1,
		Results: 5,
	}
	posts := c.repository.Fetch(q)
	qp := q.Previous()
	qp.Results = 1
	prevPosts := c.repository.Fetch(qp)
	c.viewSet.PostListView(posts, nil, prevPosts, q).Render(w)
}

func (c *PostController) GetSingle(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	key := fmt.Sprintf("%s/%s/%s", params["year"], params["month"], params["name"])
	p := c.repository.FetchOne(key)
	c.viewSet.PostSingleView(p).Render(w)
}

func (c *PostController) GetList(w http.ResponseWriter, req *http.Request) {
	s, err := strconv.Atoi(req.URL.Query().Get("start"))
	if err != nil {
		panic(err)
	}
	r, err := strconv.Atoi(req.URL.Query().Get("results"))
	if err != nil {
		panic(err)
	}
	q := &repository.PostQuery{
		Start:   uint64(s),
		Results: uint64(r),
	}
	ps := c.repository.Fetch(q)
	nextPosts := c.repository.Fetch(q.Next())
	prevPosts := c.repository.Fetch(q.Previous())
	c.viewSet.PostListView(ps, nextPosts, prevPosts, q).Render(w)
}

func (c *PostController) GetEditor(w http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	var p *entity.Post
	if key == "" {
		p = &entity.Post{}
	} else {
		p = c.repository.FetchOne(key)
		if p == nil {
			http.NotFound(w, req)
			return
		}
	}
	c.viewSet.PostEditorView(p).Render(w)
}

func (c *PostController) SubmitPost(w http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	if err := req.ParseForm(); err != nil {
		panic(err)
	}
	c.repository.Commit(&entity.Post{
		Key:     key,
		Title:   req.Form.Get("title"),
		Body:    req.Form.Get("body"),
		UrlName: req.Form.Get("url-name"),
	})
	committed := c.repository.FetchOne(key)
	http.Redirect(w, req, c.urlBuilder.LinkToPostPage(committed), http.StatusSeeOther)
}

func (c *PostController) EditPost(w http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	if err := req.ParseForm(); err != nil {
		panic(err)
	}
	c.repository.Commit(&entity.Post{
		Key:     key,
		Title:   req.Form.Get("title"),
		Body:    req.Form.Get("body"),
		UrlName: req.Form.Get("url-name"),
	})
	committed := c.repository.FetchOne(req.Form.Get("key"))
	http.Redirect(w, req, c.urlBuilder.LinkToPostPage(committed), http.StatusSeeOther)
}
