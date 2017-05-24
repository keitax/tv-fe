package controller

import (
	"net/http"

	"github.com/keitax/textvid"
	"github.com/keitax/textvid/view"
)

type AdminController struct {
	repository *textvid.Repository
	viewSet    *view.ViewSet
	config     *textvid.Config
}

func NewAdminController(r *textvid.Repository, vs *view.ViewSet, c *textvid.Config) *AdminController {
	return &AdminController{
		repository: r,
		viewSet:    vs,
		config:     c,
	}
}

func (ac *AdminController) GetIndex(w http.ResponseWriter, r *http.Request) {
	ps := ac.repository.Fetch(&textvid.PostQuery{
		Start:   1,
		Results: 0,
	})
	ac.viewSet.AdminView(ps).Render(w)
}
