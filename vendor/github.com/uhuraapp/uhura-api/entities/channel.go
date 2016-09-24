package entities

import "time"

type Channel struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"image_url"`
	Id          int64     `json:"raw_id"`
	Uri         string    `json:"id"`
	Subscribed  bool      `json:"subscribed"`
	Copyright   string    `json:"copyright"`
	Episodes    []string  `json:"episodes" sql:"-"`
	UpdatedAt   time.Time `json:"updated_at"`
	Enabled     bool      `json:"enabled"`
	Url         string    `json:"url"`
	Body        []byte    `json:"-"`
}
