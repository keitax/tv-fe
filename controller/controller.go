package controller

import (
	"net/http"

	"github.com/Sirupsen/logrus"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/dao"
	"github.com/keitax/textvid/view"
)

type Controller interface {
	GetIndex(w http.ResponseWriter, req *http.Request)
	GetSingle(w http.ResponseWriter, req *http.Request)
	GetList(w http.ResponseWriter, req *http.Request)
}

type ControllerImpl struct {
	postDao dao.PostDao
	view    view.View
	config  *config.Config
}

func New(postDao dao.PostDao, view_ view.View, config_ *config.Config) Controller {
	return &ControllerImpl{
		postDao,
		view_,
		config_,
	}
}

func (c *ControllerImpl) GetIndex(w http.ResponseWriter, req *http.Request) {
	if err := c.view.RenderIndex(w); err != nil {
		if err := c.view.Render500(w); err != nil {
			logrus.Error(err)
		}
	}
}

func (c *ControllerImpl) GetSingle(w http.ResponseWriter, req *http.Request) {
}

func (c *ControllerImpl) GetList(w http.ResponseWriter, req *http.Request) {
}
