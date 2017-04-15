package controller

import (
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/view"
)

type Controller interface {
	GetIndex(w http.ResponseWriter, req *http.Request)
	GetSingle(w http.ResponseWriter, req *http.Request)
	GetList(w http.ResponseWriter, req *http.Request)
}

type ControllerImpl struct {
	view   view.View
	config *config.Config
}

func New(view_ view.View, config_ *config.Config) Controller {
	return &ControllerImpl{view_, config_}
}

func (h *ControllerImpl) GetIndex(w http.ResponseWriter, req *http.Request) {
	if err := h.view.RenderIndex(w); err != nil {
		if err := h.view.Render500(w); err != nil {
			log.Error(err)
		}
	}
}
