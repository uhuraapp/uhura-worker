package channels

import (
	"time"

	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	"bitbucket.org/dukex/uhura-api/parser"
	"github.com/jinzhu/gorm"
)

func Create(database gorm.DB, url string) (*entities.Channel, bool) {
	channelURL, _ := helpers.ParseURL(url)
	channels, err := parser.URL(channelURL)

	if err != nil {
		log.Debug("error: %s", err)
		return nil, false
	}

	if len(channels) < 1 {
		log.Debug("error: no channel found")
		return nil, false
	}

	var ok bool
	var channel entities.Channel

	log.Debug("no error found")
	channelF := channels[0] // TODO: fix it

	log.Debug("channel UhuraID: %s", channelF.UhuraID)
	if channelF.UhuraID != "" {
		ok = database.Table(models.Channel{}.TableName()).Where("uri = ?", channelF.UhuraID).
			First(&channel).Error != gorm.RecordNotFound
	} else {
		model := translateChannel(channelF)
		log.Debug("channel: %s", model)

		ok = database.Table(models.Channel{}.TableName()).Save(&model).Error == nil
		log.Debug("is ok: %s", ok)

		if ok {
			database.Table(models.Channel{}.TableName()).First(&channel, model.Id)
			go CreateLinks(channelF.Links, channel.Id, database)
		}
	}

	return &channel, ok
}

func TranslateFromFeed(model models.Channel, channel *parser.Channel) models.Channel {
	model.Title = channel.Title
	model.Description = channel.Description
	model.Copyright = channel.Copyright
	model.ImageUrl = channel.Image
	model.Uri = helpers.MakeUri(channel.Title)
	model.Language = channel.Language
	model.UpdatedAt = time.Now()
	model.LastBuildDate = channel.LastBuildDate
	model.Url = channel.URL
	return model
}

func translateChannel(channel *parser.Channel) models.Channel {
	model := models.Channel{}
	model.CreatedAt = time.Now()
	return TranslateFromFeed(model, channel)
}

func CreateLinks(links []string, channelId int64, database gorm.DB) {
	for i := 0; i < len(links); i++ {
		u := models.ChannelURL{}
		database.Table(models.ChannelURL{}.TableName()).
			FirstOrCreate(&u, models.ChannelURL{
			ChannelId: channelId,
			Url:       links[i],
		})
	}
}
