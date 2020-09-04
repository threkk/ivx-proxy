package app

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

func (a *app) handleDownload() http.HandlerFunc {
	reDownload := regexp.MustCompile(`\$\('\.downloadlink'\)\.load\('(.+)'\)`)
	reFile := regexp.MustCompile(`downloadFollow\(event,'(.+)'\);`)

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if !isValidURL(vars["url"]) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(errorRes{
				Code:        422,
				Status:      "Unprocessable Entity",
				Description: "The parameter provided is not a URL",
			})
			return
		}

		res, err := http.Get(vars["url"])
		if err != nil || res.StatusCode != 200 {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(errorRes{
				Code:        400,
				Status:      "Bad Request",
				Description: "Error requesting the first remote resource",
			})
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(errorRes{
				Code:        400,
				Status:      "Bad Request",
				Description: "Error requesting the first remote resource",
			})
			return
		}

		links := reDownload.FindSubmatch(body)
		if len(links) != 2 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(errorRes{
				Code:        400,
				Status:      "Bad Request",
				Description: "Error requesting the second remote resource",
			})
			return
		}

		res2, err := http.Get("https://www.ivoox.com/" + string(links[1]))
		if err != nil || res.StatusCode != 200 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(errorRes{
				Code:        400,
				Status:      "Bad Request",
				Description: "Error requesting the second remote resource",
			})
			return
		}
		defer res2.Body.Close()

		body2, err := ioutil.ReadAll(res2.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(errorRes{
				Code:        400,
				Status:      "Bad Request",
				Description: "Error requesting the remote resource",
			})
			return
		}

		files := reFile.FindSubmatch(body2)
		if len(files) != 2 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(errorRes{
				Code:        400,
				Status:      "Bad Request",
				Description: "Error requesting the remote resource",
			})
			return
		}
		fmt.Println(string(files[1]))

		res3, err := http.Get(string(files[1]))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(errorRes{
				Code:        400,
				Status:      "Bad Request",
				Description: "Error requesting the remote resource",
			})
			return
		}

		defer res3.Body.Close()
		w.Header().Add("Content-Type", "audio/mp3")
		io.Copy(w, res3.Body)
		// http.Redirect(w, r, string(files[1]), http.StatusFound)
	}
}
