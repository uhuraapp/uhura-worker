package models

import (
	"time"

	"bitbucket.org/dukex/uhura-api/entities"
	"bitbucket.org/dukex/uhura-api/helpers"
	"github.com/jinzhu/gorm"
)

type Episode struct {
	Id            int64
	ChannelId     int64
	Key           string `sql:"unique"`
	SourceUrl     string `sql:"not null;unique"`
	Title         string
	Description   string
	PublishedAt   time.Time `sql:"not null"`
	Duration      string
	Uri           string
	ContentLength int64
	Type          string
	ContentType   string
	helpers.Uriable
}

func (e Episode) TableName() string {
	return "items"
}

func (e Episode) CountByChannel(database gorm.DB, channelId int64) int64 {
	var episodesCount int64

	database.Table(e.TableName()).Where("channel_id = ?", channelId).Count(&episodesCount)

	return episodesCount
}

func mapEpisodeIds(episodes []*entities.Episode) []int64 {
	ids := make([]int64, 0)
	for _, e := range episodes {
		ids = append(ids, e.Id)
	}
	return ids
}

func SetListenAttributesToEpisode(DB gorm.DB, userId int, episodes []*entities.Episode, channelURI string) {
	ids := mapEpisodeIds(episodes)

	var listeneds []*Listened

	if len(ids) > 0 {
		DB.Table(Listened{}.TableName()).
			Where("item_id IN (?)", ids).
			Where("user_id = ?", userId).
			Find(&listeneds)
	}

	mapListened := make(map[int64]*Listened, 0)
	for _, listen := range listeneds {
		mapListened[listen.ItemId] = listen
	}

	for _, episode := range episodes {
		episode.ChannelUri = channelURI

		if mapListened[episode.Id] != nil {
			episode.Listened = mapListened[episode.Id].Viewed
			episode.StoppedAt = &mapListened[episode.Id].StoppedAt
			if episode.Listened {
				z := int64(0)
				episode.StoppedAt = &z
			}
		}
	}
}
