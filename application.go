package textvid

import (
	"net/http"
	"runtime/debug"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func PanicHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.Errorf("%s: %s", err, debug.Stack())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()
	next(w, r)
}

func NewApplication(config *Config) (http.Handler, error) {
	ub := NewUrlBuilder(config)
	vs := NewViewSet(ub, config)
	re := NewRepository(config.LocalGitRepository, config.RemoteGitRepository)
	pc := NewPostController(re, vs, ub, config)
	ac := NewAdminController(re, vs, config)

	r := mux.NewRouter()
	r.HandleFunc("/", pc.GetIndex).Methods(http.MethodGet)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticDir)))).Methods(http.MethodGet)
	r.HandleFunc("/posts/", pc.GetList).Methods(http.MethodGet)
	r.HandleFunc("/posts/", pc.SubmitPost).Methods(http.MethodPost)
	r.HandleFunc("/posts/new", pc.GetEditor).Methods(http.MethodGet)
	r.HandleFunc("/posts/{id:[0-9]+}", pc.EditPost).Methods(http.MethodPost, http.MethodPut)
	r.HandleFunc("/posts/{id:[0-9]+}/edit", pc.GetEditor).Methods(http.MethodGet)
	r.HandleFunc("/{year:[0-9]{4}}/{month:0[1-9]|1[0-2]}/{name}.html", pc.GetSingle).Methods(http.MethodGet)
	r.HandleFunc("/admin", ac.GetIndex).Methods(http.MethodGet)

	app := negroni.New()
	app.UseHandler(r)
	app.UseFunc(PanicHandler)
	return app, nil
}
