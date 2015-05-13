package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURLValid(t *testing.T) {
	url, err := ParseURL("http://httpbin.org/status/200")

	assert.Nil(t, err)
	assert.Equal(t, url.Path, "/status/200", "correct path")
}

func TestParseURLSlashs(t *testing.T) {
	url, err := ParseURL("/http://httpbin.org/status/200")

	assert.Nil(t, err)
	assert.Equal(t, url.Path, "/status/200", "correct path")

	url, err = ParseURL("/////http://httpbin.org/status/200")
	assert.Nil(t, err)
	assert.Equal(t, url.Path, "/status/200", "correct path")
}

func TestParseURLNoScheme(t *testing.T) {
	url, err := ParseURL("httpbin.org/status/200")

	assert.Nil(t, err)
	assert.Equal(t, url.Path, "/status/200", "correct path")
	assert.Equal(t, url.Scheme, "http", "correct scheme")
}

func TestParseURLWorksWithHTTPS(t *testing.T) {
	url, err := ParseURL("https://httpbin.org/status/200")

	assert.Nil(t, err)
	assert.Equal(t, url.Path, "/status/200", "correct path")
	assert.Equal(t, url.Scheme, "https", "correct scheme")
}
