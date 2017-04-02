package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/handler"
	"github.com/keitax/textvid/view"
)

func main() {
	c := &config.Config{}
	v := view.New(c)
	h := handler.New(v, c)

	router := mux.NewRouter()
	router.HandleFunc("/", h.Index)

	log.Info("Listen on 8080.")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
