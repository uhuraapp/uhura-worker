package database

import (
	"log"
	"os"

	"bitbucket.org/dukex/uhura-api/models"
	"github.com/jinzhu/gorm"
	pq "github.com/lib/pq"
)

func NewPostgresql() gorm.DB {
	var database gorm.DB
	var err error

	databaseUrl, _ := pq.ParseURL(os.Getenv("DATABASE_URL"))
	database, err = gorm.Open("postgres", databaseUrl)

	if err != nil {
		log.Fatalln(err.Error())
	}

	err = database.DB().Ping()
	if err != nil {
		log.Fatalln(err.Error())
	}

	database.LogMode(os.Getenv("DEBUG") == "true")

	database.DB().Ping()
	database.DB().SetMaxIdleConns(10)
	database.DB().SetMaxOpenConns(20)

	Migrations(database)

	return database
}

func Migrations(database gorm.DB) {
	database.AutoMigrate(&models.Episode{})
	database.AutoMigrate(&models.Listened{})
	database.AutoMigrate(&models.Channel{})
	database.AutoMigrate(&models.Subscription{})
	database.AutoMigrate(&models.User{})
	database.AutoMigrate(&models.Categoriable{})
	database.AutoMigrate(&models.Category{})
	database.AutoMigrate(&models.ChannelURL{})
	database.AutoMigrate(&models.Profile{})

	database.Model(&models.Channel{}).AddIndex("idx_channel_uri", "uri")
	database.Model(&models.Channel{}).AddIndex("idx_channel_url", "url")

	database.Model(&models.Episode{}).AddIndex("idx_episode_channel_id", "channel_id")
	database.Model(&models.Episode{}).AddIndex("idx_episode_channel_id_with_published_at", "channel_id", "published_at")

	database.Model(&models.Listened{}).AddIndex("idx_listened", "item_id", "viewed", "user_id")
	database.Model(&models.Listened{}).AddIndex("idx_listened_by_channel", "channel_id", "user_id")

	database.Model(&models.Subscription{}).AddIndex("idx_subscription", "user_id")
	database.Model(&models.Subscription{}).AddIndex("idx_subscription_by_channel", "user_id", "channel_id")

	database.Model(&models.Categoriable{}).AddIndex("idx_categoriable", "channel_id", "category_id")

	database.Model(&models.User{}).AddIndex("idx_user_by_token", "api_token")
	database.Model(&models.User{}).AddIndex("idx_user_by_email", "email")
	database.Model(&models.User{}).AddUniqueIndex("idx_user_email", "email")
	database.Model(&models.User{}).AddUniqueIndex("idx_user_token", "api_token")

	database.Model(&models.ChannelURL{}).AddUniqueIndex("idx_channel_url_url", "url")

	database.Model(&models.Profile{}).AddUniqueIndex("idx_profile_key", "key")
	database.Model(&models.Profile{}).AddUniqueIndex("idx_profile_user_id", "user_id")
	database.Model(&models.Profile{}).AddIndex("idx_profile_by_key", "key")
}
