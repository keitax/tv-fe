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

	app := negroni.New()
	app.UseHandler(r)
	app.UseFunc(PanicHandler)
	return app, nil
}
