package models

import (
  "time"
)

type ChannelURL struct {
  Id        int64
  ChannelId int64  `sql:"index"`
  Url       string `sql:"type:varchar(200);unique_index"`
  CreatedAt time.Time
  UpdatedAt time.Time
}

func (c ChannelURL) TableName() string {
  return "channel_urls"
}
