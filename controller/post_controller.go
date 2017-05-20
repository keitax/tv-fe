package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/dao"
	"github.com/keitax/textvid/entity"
	"github.com/keitax/textvid/repository"
	"github.com/keitax/textvid/util"
	"github.com/keitax/textvid/view"
)

type PostController struct {
	postDao    dao.PostDao
	repository *repository.Repository
	viewSet    *view.ViewSet
	urlBuilder *util.UrlBuilder
	config     *config.Config
}

func NewPostController(postDao dao.PostDao, r *repository.Repository, vs *view.ViewSet, ub *util.UrlBuilder, config_ *config.Config) *PostController {
	return &PostController{
		postDao,
		r,
		vs,
		ub,
		config_,
	}
}

func (c *PostController) GetIndex(w http.ResponseWriter, req *http.Request) {
	q := &dao.PostQuery{
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
	q := &dao.PostQuery{
		Start:   uint64(s),
		Results: uint64(r),
	}
	ps := c.repository.Fetch(q)
	nextPosts := c.repository.Fetch(q.Next())
	prevPosts := c.repository.Fetch(q.Previous())
	c.viewSet.PostListView(ps, nextPosts, prevPosts, q).Render(w)
}

func (c *PostController) GetEditor(w http.ResponseWriter, req *http.Request) {
	ids := mux.Vars(req)["id"]
	var p *entity.Post
	if ids == "" {
		p = &entity.Post{}
	} else {
		id, err := strconv.Atoi(ids)
		if err != nil {
			panic(err)
		}
		p = c.postDao.SelectOne(int64(id))
		if p == nil {
			http.NotFound(w, req)
			return
		}
	}
	c.viewSet.PostEditorView(p).Render(w)
}

func (c *PostController) SubmitPost(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		panic(err)
	}
	p := &entity.Post{
		Title:   req.Form.Get("title"),
		Body:    req.Form.Get("body"),
		UrlName: req.Form.Get("url-name"),
	}
	c.postDao.Insert(p)
	p = c.postDao.SelectOne(int64(p.Id))
	http.Redirect(w, req, c.urlBuilder.LinkToPostPage(p), http.StatusSeeOther)
}

func (c *PostController) EditPost(w http.ResponseWriter, req *http.Request) {
	vs := mux.Vars(req)
	id, err := strconv.Atoi(vs["id"])
	if err != nil {
		panic(err)
	}
	if err := req.ParseForm(); err != nil {
		panic(err)
	}
	p := &entity.Post{
		Id:      int64(id),
		Title:   req.Form.Get("title"),
		Body:    req.Form.Get("body"),
		UrlName: req.Form.Get("url-name"),
	}
	c.postDao.Update(p)
	p = c.postDao.SelectOne(int64(id))
	http.Redirect(w, req, c.urlBuilder.LinkToPostPage(p), http.StatusSeeOther)
}
