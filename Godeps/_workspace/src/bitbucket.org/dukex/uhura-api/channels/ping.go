package channels

import (
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/jinzhu/gorm"
)

func Ping(database gorm.DB, id int64) {
	var channel models.Channel

	ok := database.Table(models.Channel{}.TableName()).
		First(&channel, id).Error != gorm.RecordNotFound

	if !ok {
		return
	}
}
