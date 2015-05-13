package helpers

import (
	"net/url"
	"regexp"
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
	return url.Parse(u)
}
