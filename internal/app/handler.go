package app

import (
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

type errorResponse struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code" json:"code"`
	Status  int      `xml:"status,attr" json:"status"`
	Message string   `xml:"message" json:"message"`
}
