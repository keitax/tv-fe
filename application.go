package textvid

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
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

func RequestLoggingHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	res := w.(negroni.ResponseWriter)
	next(w, r)
	logrus.Info(fmt.Sprintf("%s %s %d", r.Method, r.RequestURI, res.Status()))
}

func NewApplication(config *Config) (http.Handler, error) {
	ub := NewUrlBuilder(config)
	vs := NewViewSet(ub, config)
	re, err := OpenRepository(config.LocalGitRepository, config.RemoteGitRepository)
	if err != nil {
		return nil, err
	}
	re.UpdateCache()
	c := cron.New()
	c.AddFunc("*/30 * * * * *", func() {
		re.SynchronizeRemote()
	})
	c.Start()

	pc := NewPostController(re, vs, ub, config)
	ac := NewAdminController(re, vs, config)

	r := mux.NewRouter()
	r.HandleFunc("/", pc.GetIndex).Methods(http.MethodGet)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticDir)))).Methods(http.MethodGet)
	r.HandleFunc("/posts/", pc.GetList).Methods(http.MethodGet)
	r.HandleFunc("/{year:[0-9]{4}}/{month:0[1-9]|1[0-2]}/{name}.html", pc.GetSingle).Methods(http.MethodGet)
	r.HandleFunc("/admin", ac.GetIndex).Methods(http.MethodGet)

	app := negroni.New()
	app.UseHandler(r)
	app.UseFunc(PanicHandler)
	app.UseFunc(RequestLoggingHandler)
	return app, nil
}
