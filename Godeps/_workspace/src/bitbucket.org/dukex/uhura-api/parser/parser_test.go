package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestURL(t *testing.T) {
	body, err := RequestURL("http://httpbin.org/status/200")
	assert.Equal(t, body, []byte{}, "should be OK")
	assert.Nil(t, err)

	body, err = RequestURL("http://httpbin.org/status/404")
	assert.Equal(t, body, []byte(nil), "should be OK")
	assert.Equal(t, err != nil, true, "should not return a error")
}

// func TestGetChannelImageURL(t *testing.T) {
// 	channel := &rss.Channel{}
// 	url := GetChannelImageURL(channel)
// 	assert.Equal(t, url, "", "should be empty string")
//
// 	channel.Image = rss.Image{Url: "http://x"}
// 	url = GetChannelImageURL(channel)
// 	assert.Equal(t, url, "http://x", "should be equal channel image")
//
// 	channel.Extensions = make(map[string]map[string][]rss.Extension)
// 	channel.Extensions[ITUNES_EXT]["image"] = []rss.Extension{{
// 		Attrs: map[string]string{
// 			"href": "http://y",
// 		}},
// 	}
//
// 	url = GetChannelImageURL(channel)
// 	assert.Equal(t, url, "http://y", "should be equal itunes image")
// }
