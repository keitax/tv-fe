package application

import (
	"net/http"
	"runtime/debug"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/gorilla/mux"
	"github.com/keitax/textvid/config"
	"github.com/keitax/textvid/controller"
	"github.com/keitax/textvid/dao"
	"github.com/keitax/textvid/repository"
	"github.com/keitax/textvid/util"
	"github.com/keitax/textvid/view"
)

type application struct {
	router http.Handler
}

func (a *application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.Errorf("%s: %s", err, debug.Stack())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()
	a.router.ServeHTTP(w, r)
}

func New(config *config.Config) (http.Handler, error) {
	conn, err := dbr.Open("mysql", config.DataSourceName, nil)
	if err != nil {
		return nil, err
	}

	d := dao.NewPostDao(conn, config)

	ub := util.NewUrlBuilder(config)
	vs := view.NewViewSet(ub, config)
	re := repository.New(config.LocalGitRepository, config.RemoteGitRepository)
	pc := controller.NewPostController(re, vs, ub, config)
	ac := controller.NewAdminController(d, vs, config)

	r := mux.NewRouter()
	r.HandleFunc("/", pc.GetIndex)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticDir))))
	r.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			pc.SubmitPost(w, r)
			return
		}
		pc.GetList(w, r)
	})
	r.HandleFunc("/posts/new", pc.GetEditor)
	r.HandleFunc("/posts/{id:[0-9]+}", pc.EditPost)
	r.HandleFunc("/posts/{id:[0-9]+}/edit", pc.GetEditor)
	r.HandleFunc("/{year:[0-9]{4}}/{month:0[1-9]|1[0-2]}/{name}.html", pc.GetSingle)
	r.HandleFunc("/admin", ac.GetIndex)

	return &application{r}, nil
}
