package controller

import (
	"net/http"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/dao"
	"github.com/keitax/textvid/repository"
	"github.com/keitax/textvid/view"
)

type AdminController struct {
	repository *repository.Repository
	viewSet    *view.ViewSet
	config     *config.Config
}

func NewAdminController(r *repository.Repository, vs *view.ViewSet, c *config.Config) *AdminController {
	return &AdminController{
		repository: r,
		viewSet:    vs,
		config:     c,
	}
}

func (ac *AdminController) GetIndex(w http.ResponseWriter, r *http.Request) {
	ps := ac.repository.Fetch(&dao.PostQuery{
		Start:   1,
		Results: 0,
	})
	ac.viewSet.AdminView(ps).Render(w)
}
