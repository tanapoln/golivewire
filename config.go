package golivewire

import (
	"errors"
	"net/url"
	"strings"
)

var (
	baseURL         string
	DevelopmentMode bool
)

func SetBaseURL(raw string) {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	if !u.IsAbs() {
		panic(errors.New("base url is expecting full URL"))
	}
	u.RawQuery = ""
	u.Fragment = ""
	u.Path = strings.TrimSuffix(u.Path, "/")
	baseURL = u.String()
}
