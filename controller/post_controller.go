package controller

import (
	"net/http"

	"github.com/Sirupsen/logrus"

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
	if err := c.view.RenderIndex(w); err != nil {
		if err := c.view.Render500(w); err != nil {
			logrus.Error(err)
		}
	}
}

func (c *PostController) GetSingle(w http.ResponseWriter, req *http.Request) {
}

func (c *PostController) GetList(w http.ResponseWriter, req *http.Request) {
}
