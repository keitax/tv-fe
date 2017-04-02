package handler

import (
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/view"
)

type Handler interface {
	Index(http http.ResponseWriter, req *http.Request)
}

type HandlerImpl struct {
	view   view.View
	config *config.Config
}

func New(view_ view.View, config_ *config.Config) Handler {
	return &HandlerImpl{view_, config_}
}

func (h *HandlerImpl) Index(w http.ResponseWriter, req *http.Request) {
	if err := h.view.RenderIndex(w); err != nil {
		if err := h.view.Render500(w); err != nil {
			log.Error(err)
		}
	}
}
