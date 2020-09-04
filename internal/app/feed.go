package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	output "github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	input "github.com/mmcdole/gofeed"
)

func (a *app) handleRSS() http.HandlerFunc {
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
		}

		feed, err := NewFeed(vars["url"], a.baseURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(errorRes{
				Code:        400,
				Status:      "Bad Request",
				Description: "The URL could not be retrieved",
			})
		}

		rss, err := feed.Output.ToRss()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(errorRes{
				Code:        400,
				Status:      "Bad Request",
				Description: "Error decoding the RSS",
			})
		}

		if etag := r.Header.Get("ETag"); etag != "" {
			w.Header().Add("Etag", etag)
		}

		if cache := r.Header.Get("Cache-Control"); cache != "" {
			w.Header().Add("Cache-Control", cache)
		}

		w.Header().Add("Content-Type", "text/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, rss)
	}
}

type feed struct {
	url    string      // We need to keep it for the caching
	parsed *input.Feed // Maybe this can be removed?
	Output *output.Feed
}

func (f *feed) populateOutput(baseURL string) {
	f.Output = &output.Feed{
		Title: f.parsed.Title,
		Link: &output.Link{
			Href: f.parsed.Link,
		},
		Description: f.parsed.Description,
		Copyright:   f.parsed.Copyright,
		Image: &output.Image{
			Url:   f.parsed.Image.URL,
			Title: f.parsed.Image.Title,
		},
		Items: make([]*output.Item, len(f.parsed.Items)),
	}

	if f.parsed.Author != nil {
		f.Output.Author = &output.Author{
			Name:  f.parsed.Author.Email,
			Email: f.parsed.Author.Email,
		}
	}

	if f.parsed.PublishedParsed != nil {
		f.Output.Created = *f.parsed.PublishedParsed
	}

	if f.parsed.UpdatedParsed != nil {
		f.Output.Updated = *f.parsed.UpdatedParsed
	}

	for i, item := range f.parsed.Items {

		f.Output.Items[i] = &output.Item{
			Title:       item.Title,
			Link:        &output.Link{Href: item.Link},
			Source:      nil, //      &output.Link{Href: item.Link},
			Description: item.Description,
			Id:          item.GUID,
			Content:     item.Content,
		}

		if item.Author != nil {
			f.Output.Items[i].Author =
				&output.Author{
					Name:  item.Author.Name,
					Email: item.Author.Email}
		}

		if item.PublishedParsed != nil {
			f.Output.Items[i].Created = *item.PublishedParsed
		}

		if item.UpdatedParsed != nil {
			f.Output.Items[i].Updated = *item.UpdatedParsed
		}

		if len(item.Enclosures) > 0 {
			param := url.Values{
				"download": []string{item.Link}, // []string{item.Enclosures[0].URL},
			}
			f.Output.Items[i].Enclosure = &output.Enclosure{
				Url:    baseURL + "?" + param.Encode(),
				Length: item.Enclosures[0].Length,
				Type:   item.Enclosures[0].Type,
			}
		}
	}
}

func NewFeed(url string, baseURL string) (*feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fp := input.NewParser()
	fd, err := fp.ParseURLWithContext(url, ctx)

	if err != nil {
		return nil, err
	}

	f := &feed{
		url:    url,
		parsed: fd,
	}
	f.populateOutput(baseURL)
	return f, nil
}
