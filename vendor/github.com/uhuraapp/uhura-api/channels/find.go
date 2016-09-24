package channels

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/entities"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/uhuraapp/uhura-api/parser"
)

func Find(database *gorm.DB, uidORurl string) (channel entities.Channel,
	episodes entities.Episodes,
	feed *parser.Channel,
	ok bool) {
	// By default uidORurl the feedURL, hat is the first guess is this arguments
	// is a URL
	info("Find started")

	var feedURL = uidORurl

	info("Finding: " + uidORurl)

	// Try to find uidORurl in database
	err := database.Table(models.Channel{}.TableName()).
		Where("uri = ?", prepareURI(uidORurl)).
		First(&channel).Error

	if err != gorm.ErrRecordNotFound {
		// If Find the feed URL is the channel URL
		feedURL = channel.Url
		info("+ channel found")
	}

	info("Parsing url: " + feedURL)
	url, urlErr := helpers.ParseURL(feedURL)

	if urlErr != nil {
		errr(urlErr, "url error")
		return
	}

	// Find Channel with this URL
	if time.Now().Sub(channel.UpdatedAt) < (time.Hour*5) && len(channel.Body) > 5 {
		info("+ channel body is valid")
		// Use the body is valid
		feed, err = parser.Body(channel.Body, channel.Url)
	} else {
		info("+ channel body is invalid")
		// Find xml with URL
		var body []byte
		feed, body, err = parser.URL(url)
		go Save(database, feed, body)
	}

	if err != nil {
		errr(err, "Find finished")
		return channel, episodes, feed, false
	}

	info("+ translate feed to channel")
	channelFeed := ChannelEntityFromFeed(feed)
	episodes, ids := EpisodesEntityFromFeed(feed)
	channelFeed.Episodes = ids

	channel.Uri = channelFeed.Uri
	database.Table(models.Channel{}.TableName()).
		Where("id = ?", channel.Id).
		UpdateColumns(models.Channel{
			Title: channelFeed.Title,
			Uri:   channelFeed.Uri,
		})

	channelFeed.Id = channel.Id
	info("Find finished")

	return channelFeed, episodes, feed, true
}

func prepareURI(uri string) string {
	return strings.Replace(uri, "/", "", 1)
}
