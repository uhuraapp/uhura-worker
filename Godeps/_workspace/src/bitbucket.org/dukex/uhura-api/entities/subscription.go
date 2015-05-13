package entities

import (
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/jinzhu/gorm"
)

type Subscription struct {
	Id         int64  `json:"raw_id"`
	Uri        string `json:"id"`
	Title      string `json:"title"`
	ImageUrl   string `json:"image_url"`
	ToView     int64  `json:"to_view"`
	Subscribed bool   `json:"subscribed"`
}

func (s *Subscription) GetToView(database gorm.DB, userId string) int64 {
	var (
		listenedCount int64
	)

	database.Table(models.Listened{}.TableName()).
		Where("channel_id = ? AND user_id = ?", s.Id, userId).
		Count(&listenedCount)

	return models.Episode{}.CountByChannel(database, s.Id) - listenedCount
}
