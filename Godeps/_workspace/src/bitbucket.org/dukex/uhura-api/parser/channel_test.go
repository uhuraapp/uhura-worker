package parser

import (
	"testing"

	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/stretchr/testify/assert"
)

func TestChannelHasNewURL(t *testing.T) {
	channel := &Channel{}
	assert.False(t, channel.HasNewURL())

	channel = &Channel{Feed: &rss.Channel{
		Extensions: map[string]map[string][]rss.Extension{},
	}}

	assert.False(t, channel.HasNewURL())
	channel = &Channel{Feed: &rss.Channel{
		Extensions: map[string]map[string][]rss.Extension{
			"" + ITUNES_EXT: map[string][]rss.Extension{},
		},
	}}
	assert.False(t, channel.HasNewURL())

	channel = &Channel{Feed: &rss.Channel{
		Extensions: map[string]map[string][]rss.Extension{
			"" + ITUNES_EXT: map[string][]rss.Extension{
				"new-feed-url": []rss.Extension{},
			},
		},
	}}
	assert.False(t, channel.HasNewURL())

	channel = &Channel{Feed: &rss.Channel{
		Extensions: map[string]map[string][]rss.Extension{
			"" + ITUNES_EXT: map[string][]rss.Extension{
				"new-feed-url": []rss.Extension{
					rss.Extension{
						Value: "x",
					},
				},
			},
		},
	}}
	assert.True(t, channel.HasNewURL())
}

func TestChannelNewURL(t *testing.T) {
	channel := &Channel{}
	assert.Equal(t, "", channel.NewURL())

	channel = &Channel{Feed: &rss.Channel{
		Extensions: map[string]map[string][]rss.Extension{},
	}}
	assert.Equal(t, "", channel.NewURL())

	channel = &Channel{Feed: &rss.Channel{
		Extensions: map[string]map[string][]rss.Extension{
			"" + ITUNES_EXT: map[string][]rss.Extension{},
		},
	}}
	assert.Equal(t, "", channel.NewURL())

	channel = &Channel{Feed: &rss.Channel{
		Extensions: map[string]map[string][]rss.Extension{
			"" + ITUNES_EXT: map[string][]rss.Extension{
				"new-feed-url": []rss.Extension{},
			},
		},
	}}
	assert.Equal(t, "", channel.NewURL())

	channel = &Channel{Feed: &rss.Channel{
		Extensions: map[string]map[string][]rss.Extension{
			"" + ITUNES_EXT: map[string][]rss.Extension{
				"new-feed-url": []rss.Extension{
					rss.Extension{
						Value: "x",
					},
				},
			},
		},
	}}
	assert.Equal(t, "x", channel.NewURL())
}
