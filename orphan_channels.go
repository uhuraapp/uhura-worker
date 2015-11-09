package main

import (
	"bitbucket.org/dukex/uhura-api/database"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/jrallison/go-workers"
)

func orphanChannel(message *workers.Msg) {
	defer reporter(message)

	offset, _ := message.Args().Int64()

	p := database.NewPostgresql()
	var channels []models.Channel

	p.Table(models.Channel{}.TableName()).Limit(10).Offset(offset).Find(&channels)

	for _, channel := range channels {
		var users []models.Subscription
		p.Table(models.Subscription{}.TableName()).Where("channel_id = ?", channel.Id).Find(&users)
		if len(users) < 1 {
			var episodes []models.Episode
			p.Table(models.Episode{}.TableName()).Where("channel_id = ?", channel.Id).Find(&episodes)
			for _, e := range episodes {
				workers.Enqueue("delete-episode", "deleteEpisode", e.Id)
			}
			p.Table(models.Channel{}.TableName()).Where("id = ?", channel.Id).Delete(models.Channel{})
		}
	}

	if len(channels) == 0 {
		offset = -9
	}

	workers.Enqueue("orphan-channel", "orphanChannel", offset+9)
}
