package entities

import (
	"time"

	"bitbucket.org/dukex/uhura-api/models"

	"github.com/jinzhu/gorm"
)

type Episode struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Listened    bool      `json:"listened"`
	ChannelId   int64     `json:"raw_channel_id"`
	ChannelUri  string    `json:"channel_id"`
	SourceUrl   string    `json:"source_url"`
	PublishedAt time.Time `json:"published_at"`
	StoppedAt   *int64    `json:"stopped_at"`
	// Uri         string `json:"uri"`
	// Duration    string `json:"duration"`
	// Type        string `json:"type"`
}

func SetListenAttributesToEpisode(DB gorm.DB, userId int, episodes []*Episode, channelURI string) {
	ids := make([]int64, 0)

	for _, e := range episodes {
		ids = append(ids, e.Id)
	}

	var listeneds []*models.Listened

	if len(ids) > 0 {
		DB.Table(models.Listened{}.TableName()).
			Where("item_id IN (?)", ids).
			Where("user_id = ?", userId).
			Find(&listeneds)
	} else {
		listeneds = make([]*models.Listened, 0)
	}

	mapListened := make(map[int64]*models.Listened, 0)
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
