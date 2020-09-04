package app

import (
	"fmt"
	"net/http"
)

func (a *app) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Handle Index")
	}
}

func (a *app) handle404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 - Not Found", http.StatusNotFound)
	}
}
