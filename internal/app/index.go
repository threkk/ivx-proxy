package app

import (
	"html/template"
	"net/http"

	"github.com/threkk/ivx-proxy/web"
)

type index struct {
	Error string
	Proxy string
	Base  string
}

var indexTmpl = template.Must(template.New("index").Parse(web.IndexTmpl))

func (a *App) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		base := a.BaseURL
		if a.BaseURL == "" {
			base = r.URL.Host
		}
		a.Log(base)
		values := index{
			Base: base,
		}

		indexTmpl.Execute(w, values)
	}
}
