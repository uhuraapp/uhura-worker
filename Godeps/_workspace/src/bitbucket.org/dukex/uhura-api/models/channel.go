package models

import (
	"time"

	"bitbucket.org/dukex/uhura-api/helpers"
)

type Channel struct {
	Id            int64
	Title         string `sql:"not null;unique"`
	Description   string
	ImageUrl      string
	Copyright     string
	LastBuildDate string
	Url           string `sql:"not null;unique"`
	Uri           string
	Featured      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	VisitedAt     time.Time
	Language      string
	Link          string
	Loading       bool
	Colors        string
	Enabled       bool

	URLS []ChannelURL

	helpers.Uriable
}

func (c Channel) TableName() string {
	return "channels"
}
