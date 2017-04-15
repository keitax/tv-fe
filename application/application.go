package application

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/controller"
	"github.com/keitax/textvid/view"
)

func New(config *config.Config) http.Handler {
	c := controller.New(view.New(config), config)

	router := mux.NewRouter()
	router.HandleFunc("/", c.GetIndex)

	return router
}
