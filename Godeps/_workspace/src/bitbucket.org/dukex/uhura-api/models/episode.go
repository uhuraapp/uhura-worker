package models

import (
	"time"

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
