package application

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/controller"
	"github.com/keitax/textvid/dao"
	"github.com/keitax/textvid/view"
)

func New(config *config.Config) (http.Handler, error) {
	db, err := sql.Open("mysql", config.DataSourceName)
	if err != nil {
		return nil, err
	}

	pc := controller.NewPostController(dao.NewPostDao(db), view.New(config), config)

	router := mux.NewRouter()
	router.HandleFunc("/", pc.GetIndex)

	return router, nil
}
