package entities

import "time"

type Episode struct {
	Id          string    `json:"id"`
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

type Episodes []*Episode

func (e Episodes) Len() int {
	return len(e)
}

func (e Episodes) Less(i, j int) bool {
	return e[i].PublishedAt.After(e[j].PublishedAt)
}

func (e Episodes) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Episodes) IDs() []string {
	ids := make([]string, 0)
	for _, e := range e {
		ids = append(ids, e.Id)
	}
	return ids
}
