package app

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(app *App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if app.User == "" {
				next.ServeHTTP(w, r)
				return
			}

			user, pass, ok := r.BasicAuth()
			if !ok || user != app.User || pass != app.Password {
				app.HandleError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}
