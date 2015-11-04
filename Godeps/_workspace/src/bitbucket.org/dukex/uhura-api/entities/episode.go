package entities

import "time"

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
