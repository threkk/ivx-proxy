package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

func getURLandFindOne(url string, re *regexp.Regexp) (string, error) {
	if !isValidURL(url) {
		return "", fmt.Errorf("The parameter provided is not a URL: %s", url)
	}

	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		return "", fmt.Errorf("Error requesting the remote resource: %s", url)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading the remote resource: %s", url)
	}

	elements := re.FindSubmatch(body)
	if len(elements) != 2 {
		return "", fmt.Errorf("Element not found.")
	}

	return string(elements[1]), nil
}

func (a *app) handleDownload() http.HandlerFunc {
	reDownload := regexp.MustCompile(`\$\('\.downloadlink'\)\.load\('(.+)'\)`)
	reFile := regexp.MustCompile(`downloadFollow\(event,'(.+)'\);`)

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		link, err := getURLandFindOne(vars["url"], reDownload)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(errorRes{
				Code:        400,
				Status:      "Bad Request",
				Description: err.Error(),
			})
			return
		}

		file, err := getURLandFindOne("https://www.ivoox.com/"+link, reFile)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(errorRes{
				Code:        400,
				Status:      "Bad Request",
				Description: err.Error(),
			})
			return
		}

		http.Redirect(w, r, file, http.StatusMovedPermanently)
	}
}
