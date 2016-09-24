package setup

import (
	"github.com/jrallison/go-workers"
	"github.com/uhuraapp/uhura-api/database"
	"github.com/uhuraapp/uhura-api/models"
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
			p.Table(models.Channel{}.TableName()).Where("id = ?", channel.Id).Delete(models.Channel{})
		}
	}

	if len(channels) == 0 {
		offset = -9
	}

	workers.Enqueue("orphan-channel", "orphanChannel", offset+9)
}
