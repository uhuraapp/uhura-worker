package main

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"bitbucket.org/dukex/uhura-api/channels"
	"bitbucket.org/dukex/uhura-api/database"
	"bitbucket.org/dukex/uhura-api/models"
	"bitbucket.org/dukex/uhura-api/parser"

	"github.com/jinzhu/gorm"
	"github.com/jrallison/go-workers"
)

const ImageHost string = "http://arcane-forest-5063.herokuapp.com"

var p gorm.DB

func main() {
	redis, err := url.Parse(os.Getenv("REDIS_URL"))

	if err != nil {
		panic("REDIS_URL error, " + err.Error())
	}

	password, _ := redis.User.Password()

	workers.Configure(map[string]string{
		"server":   redis.Host,
		"password": password,
		"database": "0",
		"pool":     "20",
		"process":  "1",
	})

	workers.Process("sync", sync, 18)
	workers.Process("sync-low", sync, 2)

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	go workers.StatsServer(port)

	p = database.NewPostgresql()
	workers.Run()
}

func sync(message *workers.Msg) {
	var channel models.Channel

	id, err := message.Args().Int64()
	checkError(err)

	// 	// [x] channel := FindChannel(xml)
	err = p.Table(models.Channel{}.TableName()).First(&channel, id).Error
	checkError(err)

	channelURL, err := url.Parse(channel.Url)
	checkError(err)

	// 	// [x] xml := ParserXML(body)
	channelsFeed, errors := parser.URL(channelURL)
	if len(errors) > 0 {
		panic(errors[0])
	}

	channelFeed := channelsFeed[0]

	//  // [x] UpdateChannel(channel, xml)
	updateChannel(&channel, channelFeed)

	// 	// [x] CacheImage(channel)
	cacheImage(&channel)

	p.Save(&channel)

	channels.CreateLinks(channelFeed.Links, channel.Id, p)

	// 	// [ ] episodes := FindOrCreateEpisodes(channel, xml)
	// 	// [ ] GetDelayBetweenEpisodes(episodes)
	// 	// [ ] SetNewRun(channel)
}

func updateChannel(model *models.Channel, feed *parser.Channel) *models.Channel {
	return channels.TranslateFromFeed(model, feed)
}

func cacheImage(model *models.Channel) {
	currentImageURL := model.ImageUrl
	if strings.Contains(currentImageURL, ImageHost) {
		// Check image is OK
		return
	}

	resp, err := http.Get(ImageHost + "/resolve?url=" + currentImageURL)
	if err != nil {
		return
	}

	newImageURL := resp.Request.URL.String()

	if resp.StatusCode == 200 && strings.Contains(newImageURL, ImageHost+"/cache") {
		model.ImageUrl = newImageURL
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
