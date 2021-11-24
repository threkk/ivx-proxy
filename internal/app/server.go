package app

import (
	"encoding/xml"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// App ivoox proxy entry point.
type App struct {
	BaseURL  string
	User string
	Password string
	router   *mux.Router
	logger   *log.Logger
}

func (a *App) routes() {
	a.router.HandleFunc("/feed", a.handleRSS()).Methods("GET").Queries("url", "{url}")
	a.router.HandleFunc("/dl", a.handleDownload()).Methods("GET").Queries("url", "{url}")
	a.router.HandleFunc("/health", a.handleHealth()).Methods("GET")
	a.router.HandleFunc("/", a.handleIndex()).Methods("GET")
	a.router.PathPrefix("/").HandlerFunc(a.handle404()).Methods("GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH")
}

func (a *App) Log(format string, vars ...interface{}) {
	a.logger.Printf("LOG - "+format, vars...)
}

func (a *App) Warn(format string, vars ...interface{}) {
	a.logger.Printf("WARN - "+format, vars...)
}

func (a *App) Err(format string, vars ...interface{}) {
	a.logger.Printf("ERROR - "+format, vars...)
}

// HandleError Generates an error reponse based on the status code and a custom
// message.
func (a *App) HandleError(w http.ResponseWriter, status int, msg string) {
	output, err := xml.Marshal(errorResponse{
		Code:    http.StatusText(status),
		Status:  status,
		Message: msg,
	})

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	w.Write(output)
}

// ServeHTTP Starts the server.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// NewApp Creates a new instance of ivoox proxy
func NewApp(baseURL string, user string, password string) *App {
	a := &App{
		BaseURL:  baseURL,
		User: user,
		Password: password,
	}

	// Initialise fields
	a.logger = log.New(os.Stderr, "[app]", log.LstdFlags)
	a.router = mux.NewRouter()

	// Set routes.
	a.routes()

	// Set middleware
	a.router.Use(loggingMiddleware)
	a.router.Use(authMiddleware(a))
	a.router.Use(handlers.ProxyHeaders)
	a.router.Use(corsMiddleware)
	a.router.Use(mux.CORSMethodMiddleware(a.router))
	a.router.Use(handlers.RecoveryHandler())

	return a
}
