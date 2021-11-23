package app

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

func (a *App) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "Handle Index")
	}
}

func (a *App) handle404() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		a.HandleError(w, 404, "These are not the droids you are looking for.")
	}
}

func (a *App) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(struct {
			Status string `json:"status"`
		}{
			Status: "UP",
		})
	}
}

type errorResponse struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code" json:"code"`
	Status  int      `xml:"status,attr" json:"status"`
	Message string   `xml:"message" json:"message"`
}
