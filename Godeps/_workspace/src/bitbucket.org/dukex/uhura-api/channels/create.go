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
	channelF, err := parser.URL(channelURL)

	if err != nil {
		log.Debug("error: %s", err)
		return nil, false
	}

	var ok bool
	var channel entities.Channel

	log.Debug("no error found")

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
			// TODO: add to sync queue
		}
	}

	return &channel, ok
}

func TranslateFromFeedToEntity(entity entities.Channel, channel *parser.Channel) entities.Channel {
	entity.Uri = channel.ID
	entity.Title = channel.Title
	entity.Description = channel.Description
	entity.Copyright = channel.Copyright
	entity.ImageUrl = channel.Image
	entity.UpdatedAt = time.Now()
	return entity
}

func TranslateEpisodesFromFeedToEntity(channel *parser.Channel) ([]*entities.Episode, []int64) {
	episodes := make([]*entities.Episode, 0)
	ids := make([]int64, 0)

	for i, episode := range channel.Episodes {
		s := int64(0)
		id := int64(i) + time.Now().Unix()
		episodes = append(episodes, &entities.Episode{
			Id:          id,
			Title:       episode.Title,
			Description: episode.Description,
			SourceUrl:   episode.Source,
			StoppedAt:   &s,
			ChannelUri:  channel.ID,
		})
		ids = append(ids, id)
	}

	return episodes, ids
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
