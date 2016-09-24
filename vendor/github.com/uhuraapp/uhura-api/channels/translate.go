package channels

import (
	"sort"
	"time"

	"github.com/uhuraapp/uhura-api/entities"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/uhuraapp/uhura-api/parser"
)

func ChannelEntityFromFeed(channel *parser.Channel) entities.Channel {
	entity := entities.Channel{}
	entity.Title = channel.Title
	entity.Description = channel.Description
	entity.Copyright = channel.Copyright
	entity.ImageUrl = channel.Image
	entity.Uri = helpers.MakeUri(channel.Title)
	entity.UpdatedAt = time.Now()
	entity.Enabled = true
	return entity
}

func ChannelModelFromFeed(channel *parser.Channel) models.Channel {
	model := models.Channel{}
	model.Title = channel.Title
	model.Description = channel.Description
	model.Copyright = channel.Copyright
	model.ImageUrl = channel.Image
	model.Uri = channel.URI
	model.Language = channel.Language
	model.UpdatedAt = time.Now()
	model.LastBuildDate = channel.LastBuildDate
	model.Url = channel.URL
	model.Body = channel.Body
	return model
}

func EpisodesEntityFromFeed(channel *parser.Channel) (entities.Episodes, []string) {
	episodes := make(entities.Episodes, 0)
	ids := make([]string, 0)

	for _, episode := range channel.Episodes {
		s := int64(0)
		id := helpers.MakeUri(episode.Title)
		episodes = append(episodes, &entities.Episode{
			Id:          id,
			Title:       episode.Title,
			Description: episode.Description,
			PublishedAt: episode.PublishedAt,
			SourceUrl:   episode.Source,
			StoppedAt:   &s,
			ChannelUri:  channel.URI,
		})
	}

	sort.Sort(episodes)

	for i := range episodes {
		ids = append(ids, episodes[i].Id)
	}

	return episodes, ids
}
