package app

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

const (
	IVOOX_AUDIO_ML    = "http://www.ivoox.com/listen_mn_id_1.m4a?internal=HTML5"
	IVOOX_AUDIO_ID_RE = `(?m)var\s+idaudio\s+=\s+(\d+)`
)

var reAudio = regexp.MustCompile(IVOOX_AUDIO_ID_RE)

func (a *App) handleDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if url, ok := vars["url"]; !ok || !isValidURL(url) {
			a.HandleError(w, http.StatusUnprocessableEntity, "The parameter provided is missing or not a URL")
			return
		}

		entry, err := fetch(vars["url"])
		if err != nil {
			a.Err(err.Error())
			a.HandleError(w, http.StatusBadRequest, "The URL could not be retrieved")
			return
		}

		matches := reAudio.FindStringSubmatch(entry)
		if matches == nil {
			a.Err(fmt.Sprintf("File not found in the entry. url=%s", vars["url"]))
			a.HandleError(w, http.StatusInternalServerError, "File not found in the entry")
			return
		}

		audioId := matches[1]
		// Believe it or not, this is how it works in the original version
		file := strings.Replace(IVOOX_AUDIO_ML, "id", audioId, 1)

		http.Redirect(w, r, file, http.StatusMovedPermanently)
	}
}
