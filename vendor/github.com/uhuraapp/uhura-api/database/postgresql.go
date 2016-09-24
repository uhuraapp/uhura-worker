package database

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	pq "github.com/lib/pq"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
)

func NewPostgresql() (database *gorm.DB) {
	var err error

	databaseUrl, e := pq.ParseURL(os.Getenv("DATABASE_URL"))
	log.Println(e)
	log.Println("Connecting " + databaseUrl + "...")
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

	if os.Getenv("MIGRATIONS") == "true" {
		Migrations(database)
	}

	return database
}

func Migrations(database *gorm.DB) {
	database.DB().Exec("create extension \"uuid-ossp\";")

	database.AutoMigrate(&models.Listened{})
	database.AutoMigrate(&models.Channel{})
	database.AutoMigrate(&models.Subscription{})
	database.AutoMigrate(&models.User{})
	database.AutoMigrate(&models.Categoriable{})
	database.AutoMigrate(&models.Category{})
	database.AutoMigrate(&models.ChannelURL{})

	database.Model(&models.Channel{}).AddIndex("idx_channel_uri", "uri")
	database.Model(&models.Channel{}).AddIndex("idx_channel_url", "url")

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

	// Search
	_, err := database.DB().Exec(`DROP FUNCTION IF EXISTS channels_search_trigger() CASCADE`)
	if err != nil {
		log.Println(err)
	}

	database.Exec("DROP EXTENSION IF EXISTS unaccent")
	database.Exec("ALTER TABLE channels DROP COLUMN IF EXISTS tsv")
	database.Exec("DROP INDEX IF EXISTS channels_tsv_idx")

	database.Exec("DROP TRIGGER tsvectorupdate")

	var plays []struct {
		Id    int
		Title string
	}

	database.Table(models.Listened{}.TableName()).
		Select("user_items.id AS id, items.title AS title").
		Joins("JOIN items ON items.id = user_items.item_id").
		Find(&plays)

	for _, play := range plays {
		database.Table(models.Listened{}.TableName()).
			Where("id = ?", play.Id).
			UpdateColumn("item_uid", helpers.MakeUri(play.Title))
	}
}
