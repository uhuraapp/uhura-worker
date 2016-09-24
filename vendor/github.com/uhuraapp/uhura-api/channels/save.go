package channels

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/uhuraapp/uhura-api/parser"
)

func Save(db *gorm.DB, channel *parser.Channel, body []byte) {
	var model models.Channel

	notFound := db.Table(models.Channel{}.TableName()).
		Where("url = ?", channel.URL).
		Or("uri = ?", channel.URI).
		Limit(1).
		Find(&model).RecordNotFound()

	if notFound {
		model = ChannelModelFromFeed(channel)
		db.Save(&model)
	}

	db.Table(models.Channel{}.TableName()).
		Where("id = ?", model.Id).
		Updates(map[string]interface{}{
			"body":       body,
			"updated_at": time.Now(),
		})
}
