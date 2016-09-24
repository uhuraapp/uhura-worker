package entities

import "time"

type Category struct {
  Id          int64     `json:"raw_id"`
  Uri         string    `json:"id"`
  Name        string    `json:"name"`
  CreatedAt   time.Time `json:"created_at"`
  UpdatedAt   time.Time `json:"updated_at"`
  ChannelsIDs []string  `json:"channels" sql:"-"`
}
