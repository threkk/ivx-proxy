package app

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mmcdole/gofeed"
)

const (
	IVOOX_HOST          = "www.ivoox.com"
	IVOOX_ORIGINALS_URL = "ivoox.com/originals"
	IVOOX_RE_GENERATOR  = `(?m)<generator>iVoox<\/generator>`
	IVOOX_RE_ENCLOSURE  = `(?m)<enclosure\s+url="\S+"\s+type="audio/mpeg"\s+length="(\d+)"/>`
	ENCLOSURE_TPL       = "<enclosure url=\"%s\" type=\"audio/mpeg\" length=\"%s\"/>"
)

var reGenerator = regexp.MustCompile(IVOOX_RE_GENERATOR)
var reEnclosure = regexp.MustCompile(IVOOX_RE_ENCLOSURE)

func (a *App) handleRSS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if url, ok := vars["url"]; !ok || !isValidURL(url) {
			a.HandleError(w, http.StatusUnprocessableEntity, "The parameter provided is missing or not a URL")
			return
		}

		podcastURL, err := url.Parse(vars["url"])
		if err != nil || podcastURL.Host != IVOOX_HOST {
			a.Err(err.Error())
			a.HandleError(w, http.StatusUnprocessableEntity, "The podcast must be hosted on ivoox.com")
			return
		}

		rawRSS, err := fetch(podcastURL.String())
		if err != nil {
			a.Err(err.Error())
			a.HandleError(w, http.StatusBadRequest, "The requested url is not available at the moment")
		}

		feed, err := NewFeed(rawRSS, r)
		if err != nil {
			a.Err(err.Error())
			a.HandleError(w, http.StatusInternalServerError, "The URL could not be retrieved")
			return
		}

		if user, pass, ok := r.BasicAuth(); ok {
			feed.UserInfo = url.UserPassword(user, pass)
		}

		if !feed.IsIvooxOriginals {
			a.HandleError(w, http.StatusBadRequest, "The podcast is not an iVoox originals")
			return
		}

		if err != nil {
			a.Err(err.Error())
			a.HandleError(w, http.StatusInternalServerError, "The response could not be generated")
			return
		}

		if etag := r.Header.Get("ETag"); etag != "" {
			w.Header().Set("ETag", etag)
		}

		if cache := r.Header.Get("Cache-Control"); cache != "" {
			w.Header().Set("Cache-Control", cache)
		}

		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(feed.Patch()))
	}
}

// Feed Contains the feed url, the parsed RSS feed and the modified output.
type Feed struct {
	IsIvooxOriginals bool
	UserInfo         *url.Userinfo
	BaseURL          string
	feed             *gofeed.Feed
	rawRSS           string
	scheme           string
}

// String This is a bit trickier than expected. gofeed does not allow to export
// to XML, only to JSON. This is a problem as we want to keep the XML output but
// we don't want to reimplement the whole RSS parser.
//
// The hack consists on using the parsed version to generate the URL and later
// replace them in the right order in the original RSS. It is ugly, but it works
// and it is better than reimplmenting the library.
func (f *Feed) Patch() string {
	// Retrieve all the links per entry and generate the new ones.
	links := make([]string, len(f.feed.Items))
	for idx, item := range f.feed.Items {
		links[idx] = f.generateURL(item.Link)
	}

	// Closures. There is always only one per item. We need to find them all and
	// create new ones using the links before. Because we know that each item
	// has one closure, we can make "dangerous" index accesses.
	enclosures := make([]string, len(f.feed.Items))
	for idx, enclosure := range reEnclosure.FindAllStringSubmatch(f.rawRSS, -1) {
		enclosures[idx] = fmt.Sprintf(ENCLOSURE_TPL, links[idx], enclosure[1])
	}

	// We generate the new output by replacing every enclosure for the new one.
	idx := 0
	output := reEnclosure.ReplaceAllStringFunc(f.rawRSS, func(_ string) string {
		idx++
		return enclosures[idx-1]
	})

	// Proud iVoox Proxy generator
	output = reGenerator.ReplaceAllString(output, "iVoox Proxy")
	return output
}

// NewFeed Creates a new feed based on a string.
func NewFeed(content string, r *http.Request) (*Feed, error) {
	fp := gofeed.NewParser()
	fd, err := fp.ParseString(content)

	if err != nil {
		return nil, err
	}

	isIvooxOriginals := fd.Generator == "iVoox" && strings.Contains(content, IVOOX_ORIGINALS_URL)

	f := &Feed{
		IsIvooxOriginals: isIvooxOriginals,
		BaseURL:          r.URL.Host,
		feed:             fd,
		rawRSS:           content,
		scheme:           r.URL.Scheme,
	}

	return f, nil
}

func (f *Feed) generateURL(link string) string {
	query := url.Values{}
	query.Add("url", link)

	newURL := url.URL{
		Host:     f.BaseURL,
		Path:     "/dl",
		Scheme:   f.scheme,
		User:     f.UserInfo,
		RawQuery: query.Encode(),
	}

	return newURL.String()
}
