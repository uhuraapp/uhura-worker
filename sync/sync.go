package sync

import (
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/uhuraapp/uhura-api/channels"
	//"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
	//"github.com/uhuraapp/uhura-api/entities"
	"github.com/uhuraapp/uhura-api/parser"
)

const (
	weekHours time.Duration = 24 * 7
)

var (
	imageHosts     []string = []string{"https://images1.uhura.io", "https://images2.uhura.io"}
	imageHostRegEx          = regexp.MustCompile(`.+images\d.uhura.io`)
)

func Sync(channelID int64, p *gorm.DB) (model models.Channel, feed *parser.Channel) {
	model = GetChannel(channelID, p)

	_, _, feed, ok := channels.Find(p, model.Uri)
	if ok {
		createCategory(feed, model, p)
	}

	return model, feed
}

func GetChannel(id int64, p *gorm.DB) (model models.Channel) {
	err := p.Table(models.Channel{}.TableName()).First(&model, id).Error
	checkError(err)

	return model
}

// TODO: remote it
func cacheImage(model models.Channel) models.Channel {
	currentImageURL := model.ImageUrl

	if imageHostRegEx.MatchString(currentImageURL) {
		return model
	}

	imageHost := random(imageHosts)
	resp, err := http.Get(imageHost + "/resolve?url=" + currentImageURL)
	if err != nil {
		return model
	}

	newImageURL := resp.Request.URL.String()

	if resp.StatusCode == 200 && strings.Contains(newImageURL, imageHost+"/cache") {
		model.ImageUrl = newImageURL
	}

	return model
}

// GetNextRun returns the next run to channel
func GetNextRun(feed *parser.Channel) (time.Time, error) {
	now := time.Now()

	if len(feed.Episodes) > 1 {
		last, errLast := feed.Episodes[0].FixPubDate()
		if errLast != nil {
			return now, errLast
		}

		penultimate, errPenutimate := feed.Episodes[1].FixPubDate()
		if errPenutimate != nil {
			return now, errLast
		}

		// The next run is the duration of last less penultimate episode
		nextRunAt := last.Add(last.Sub(penultimate))

		// If next run date was a old date
		if !nextRunAt.After(now) {
			return now.Add(time.Hour * weekHours), nil
		}

		return nextRunAt, nil
	}

	return now.Add(time.Hour * weekHours), nil
}

func createCategory(feed *parser.Channel, model models.Channel, p *gorm.DB) {
	for _, data := range feed.Categories {
		var category models.Category

		p.Table(models.Category{}.TableName()).
			Where("name = ?", data.Name).
			FirstOrCreate(&category)

		p.Table(models.Categoriable{}.TableName()).
			Where("channel_id = ? AND category_id = ?", model.Id, category.Id).
			FirstOrCreate(&models.Categoriable{})
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func random(slice []string) string {
	rand.Seed(time.Now().Unix())

	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice[0]
}

func unique(s []string) []string {
	m := map[string]bool{}
	t := []string{}
	for _, v := range s {
		if _, seen := m[v]; !seen {
			t = append(t, v)
			m[v] = true
		}
	}
	return t
}
