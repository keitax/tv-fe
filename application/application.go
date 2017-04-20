package application

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/controller"
	"github.com/keitax/textvid/dao"
	"github.com/keitax/textvid/util"
	"github.com/keitax/textvid/view"
)

func New(config *config.Config) (http.Handler, error) {
	db, err := sql.Open("mysql", config.DataSourceName)
	if err != nil {
		return nil, err
	}

	v := view.New(util.NewUrlBuilder(config), config)
	pc := controller.NewPostController(dao.NewPostDao(db, config), v, config)

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticDir))))
	router.HandleFunc("/{year:[0-9]{4}}/{month:0[1-9]|1[0-2]}/{name}.html", pc.GetSingle)
	router.HandleFunc("/", pc.GetIndex)

	return router, nil
}
