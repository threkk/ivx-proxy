package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

type app struct {
	router  *mux.Router
	db      []string
	baseURL string
}

func (a *app) routes() {
	a.router.HandleFunc("/", a.handleRSS()).Methods("GET").Queries("rss", "{url}")
	a.router.HandleFunc("/", a.handleDownload()).Methods("GET").Queries("download", "{url}")
	a.router.HandleFunc("/", a.handleIndex())
	a.router.PathPrefix("/").HandlerFunc(a.handle404())
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func NewApp() *app {
	a := &app{
		baseURL: "localhost:3000",
	}
	a.router = mux.NewRouter()
	a.routes()
	return a
}
