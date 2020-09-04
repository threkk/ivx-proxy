package app

import "net/url"

func isValidURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil || u.Host == "" {
		return false
	}

	return true
}

type errorRes struct {
	Code        int    `json:"code"`
	Status      string `json:"status"`
	Description string `json:"description"`
}
