package app

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// App ivoox proxy entry point.
type App struct {
	router  *mux.Router
	baseURL string
}

func (a *App) routes() {
	a.router.HandleFunc("/", a.handleRSS()).Methods("GET").Queries("feed", "{url}")
	a.router.HandleFunc("/", a.handleDownload()).Methods("GET").Queries("dl", "{url}")
	a.router.HandleFunc("/", a.handleIndex()).Methods("GET")
	a.router.PathPrefix("/").HandlerFunc(a.handle404()).Methods("GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH")
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// NewApp Creates a new instance of ivoox proxy
func NewApp(baseURL string) *App {
	a := &App{
		baseURL: baseURL,
	}
	a.router = mux.NewRouter()
	a.routes()

	a.router.Use(func(h http.Handler) http.Handler {
		return handlers.CombinedLoggingHandler(os.Stdout, h)
	})
	a.router.Use(handlers.ProxyHeaders)
	a.router.Use(corsMiddleware)
	a.router.Use(mux.CORSMethodMiddleware(a.router))
	a.router.Use(handlers.RecoveryHandler())

	return a
}
