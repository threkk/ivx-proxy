package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func isValidURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil || u.Host == "" {
		return false
	}

	return true
}

func fetch(u string) (string, error) {
	if !isValidURL(u) {
		return "", fmt.Errorf("the parameter provided is not a URL: %s", u)
	}

	res, err := http.Get(u)
	if err != nil || res.StatusCode != 200 {
		return "", fmt.Errorf("error requesting the remote resource: %s", u)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the remote resource: %s", u)
	}

	return string(body), nil
}
