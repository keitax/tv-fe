package textvid

import (
	"net/http"
)

type AdminController struct {
	repository *Repository
	viewSet    *ViewSet
	config     *Config
}

func NewAdminController(r *Repository, vs *ViewSet, c *Config) *AdminController {
	return &AdminController{
		repository: r,
		viewSet:    vs,
		config:     c,
	}
}

func (ac *AdminController) GetIndex(w http.ResponseWriter, r *http.Request) {
	ps := ac.repository.Fetch(&PostQuery{
		Start:   1,
		Results: 0,
	})
	ac.viewSet.AdminView(ps).Render(w)
}
