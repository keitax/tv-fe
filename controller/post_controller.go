package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/dao"
	"github.com/keitax/textvid/view"
)

type PostController struct {
	postDao dao.PostDao
	view    view.View
	config  *config.Config
}

func NewPostController(postDao dao.PostDao, view_ view.View, config_ *config.Config) *PostController {
	return &PostController{
		postDao,
		view_,
		config_,
	}
}

func (c *PostController) GetIndex(w http.ResponseWriter, req *http.Request) {
	q := &dao.PostQuery{
		Start:   1,
		Results: 5,
	}
	posts, err := c.postDao.SelectByQuery(q)
	if err != nil {
		c.fatalResponse(w, err)
		return
	}
	qp := q.Previous()
	qp.Results = 1
	prevPosts, err := c.postDao.SelectByQuery(qp)
	if err != nil {
		c.fatalResponse(w, err)
		return
	}
	if err := c.view.RenderTemplate("post_list.tmpl", w, map[string]interface{}{
		"posts":         posts,
		"PreviousPosts": prevPosts,
		"CurrentQuery":  q,
	}); err != nil {
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
	posts, err := c.postDao.SelectByQuery(&dao.PostQuery{
		Start:   1,
		Results: 1,
		Year:    year,
		Month:   time.Month(month),
		UrlName: urlName,
	})
	if err != nil {
		c.fatalResponse(w, err)
		return
	}
	if len(posts) <= 0 {
		http.NotFound(w, req)
		return
	}
	p, err := c.postDao.SelectOne(posts[0].Id)
	if err != nil {
		c.fatalResponse(w, err)
		return
	}
	if err := c.view.RenderTemplate("post_single.tmpl", w, map[string]interface{}{
		"post": p,
	}); err != nil {
		c.fatalResponse(w, err)
		return
	}
}

func (c *PostController) GetList(w http.ResponseWriter, req *http.Request) {
}

func (c *PostController) fatalResponse(w http.ResponseWriter, err error) {
	logrus.Error(err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
