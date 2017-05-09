package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/dao"
	"github.com/keitax/textvid/entity"
	"github.com/keitax/textvid/util"
	"github.com/keitax/textvid/view"
)

type PostController struct {
	postDao    dao.PostDao
	viewSet    *view.ViewSet
	urlBuilder *util.UrlBuilder
	config     *config.Config
}

func NewPostController(postDao dao.PostDao, vs *view.ViewSet, ub *util.UrlBuilder, config_ *config.Config) *PostController {
	return &PostController{
		postDao,
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
	posts := c.postDao.SelectByQuery(q)
	qp := q.Previous()
	qp.Results = 1
	prevPosts := c.postDao.SelectByQuery(qp)
	if err := c.viewSet.PostListView(posts, nil, prevPosts, q).Render(w); err != nil {
		c.fatalResponse(w, err)
		return
	}
}

func (c *PostController) GetSingle(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	year, err := strconv.Atoi(params["year"])
	if err != nil {
		c.fatalResponse(w, err)
		return
	}
	month, err := strconv.Atoi(params["month"])
	if err != nil {
		c.fatalResponse(w, err)
		return
	}
	urlName := params["name"]
	posts := c.postDao.SelectByQuery(&dao.PostQuery{
		Start:   1,
		Results: 1,
		Year:    year,
		Month:   time.Month(month),
		UrlName: urlName,
	})
	if len(posts) <= 0 {
		http.NotFound(w, req)
		return
	}
	p := c.postDao.SelectOne(posts[0].Id)
	if err := c.viewSet.PostSingleView(p).Render(w); err != nil {
		c.fatalResponse(w, err)
		return
	}
}

func (c *PostController) GetList(w http.ResponseWriter, req *http.Request) {
	s, err := strconv.Atoi(req.URL.Query().Get("start"))
	if err != nil {
		c.fatalResponse(w, err)
		return
	}
	r, err := strconv.Atoi(req.URL.Query().Get("results"))
	if err != nil {
		c.fatalResponse(w, err)
		return
	}
	q := &dao.PostQuery{
		Start:   uint64(s),
		Results: uint64(r),
	}
	ps := c.postDao.SelectByQuery(q)
	nextPosts := c.postDao.SelectByQuery(q.Next())
	prevPosts := c.postDao.SelectByQuery(q.Previous())
	if err := c.viewSet.PostListView(ps, nextPosts, prevPosts, q).Render(w); err != nil {
		c.fatalResponse(w, err)
		return
	}
}

func (c *PostController) GetEditor(w http.ResponseWriter, req *http.Request) {
	vs := mux.Vars(req)
	id, err := strconv.Atoi(vs["id"])
	if err != nil {
		c.fatalResponse(w, err)
		return
	}
	p := c.postDao.SelectOne(int64(id))
	if p == nil {
		http.NotFound(w, req)
		return
	}
	if err := c.viewSet.PostEditorView(p).Render(w); err != nil {
		c.fatalResponse(w, err)
		return
	}
}

func (c *PostController) EditPost(w http.ResponseWriter, req *http.Request) {
	vs := mux.Vars(req)
	id, err := strconv.Atoi(vs["id"])
	if err != nil {
		c.fatalResponse(w, err)
		return
	}
	if err := req.ParseForm(); err != nil {
		c.fatalResponse(w, err)
		return
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

func (c *PostController) fatalResponse(w http.ResponseWriter, err error) {
	logrus.Error(err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
