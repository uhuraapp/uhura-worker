package helpers

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var hasSchemeRegx = regexp.MustCompile(`^https?:\/\/`)
var slashOnBeginRegx = regexp.MustCompile(`^\/{1,}`)

func ParseURL(u string) (*url.URL, error) {
	if slashOnBeginRegx.Match([]byte(u)) {
		u = slashOnBeginRegx.ReplaceAllString(u, "")
	}

	if !hasSchemeRegx.Match([]byte(u)) {
		u = "http://" + u
	}

	if !strings.Contains(u, ".") {
		return nil, errors.New("invalid url")
	}

	return url.Parse(u)
}
