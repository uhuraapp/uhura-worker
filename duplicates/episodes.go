package duplates

import (
	"sort"
	"strings"

	"github.com/jinzhu/gorm"

	"bitbucket.org/dukex/uhura-api/models"
)

type episode struct {
	Title string
	Count int64
	ID    int64
}

// Episodes find duplicated episodes
func Episodes(DB gorm.DB) []int64 {
	var episodes []episode

	DB.Table(models.Episode{}.TableName()).Select("items.title as title, o.dupeCount as count, items.id as id").Joins("INNER JOIN (SELECT title, channel_id, COUNT(*) as dupeCount FROM items GROUP BY title,channel_id HAVING COUNT(*) > 1) o on o.title = items.title AND o.channel_id = items.channel_id").Scan(&episodes)

	organizedEpisodes := organizeDuplicates(episodes)

	episodesToDelete := make([]int64, 0)
	for _, es := range organizedEpisodes {
		e, others := lastAndOthersEpisodes(es)
		updatePlays(e, others, DB)
		for _, other := range others {
			episodesToDelete = append(episodesToDelete, other.ID)
		}
	}

	return episodesToDelete
}

//
func organizeDuplicates(episodes []episode) map[string][]episode {
	duplicateEpisodes := make(map[string][]episode)
	for _, e := range episodes {
		duplicateEpisodes[strings.ToLower(e.Title)] = append(duplicateEpisodes[e.Title], e)
	}
	return duplicateEpisodes
}

func lastAndOthersEpisodes(episodes []episode) (episode, []episode) {
	sort.Sort(episodeByID(episodes))
	newEpisodes := episodes[:len(episodes)-1]
	return episodes[len(episodes)-1], newEpisodes
}

func updatePlays(e episode, others []episode, DB gorm.DB) {
	plays := make([]models.Listened, 0)
	for _, o := range others {
		otherPlays := getPlays(o, DB)
		plays = append(plays, otherPlays...)
	}

	for _, l := range plays {
		DB.Table(models.Listened{}.TableName()).Where("id = ?", l.Id).Update("item_id", e.ID)
	}
}

func getPlays(e episode, DB gorm.DB) []models.Listened {
	plays := make([]models.Listened, 0)
	DB.Table(models.Listened{}.TableName()).Where("item_id = ?", e.ID).Find(&plays)
	return plays
}

//

type episodeByID []episode

func (a episodeByID) Len() int           { return len(a) }
func (a episodeByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a episodeByID) Less(i, j int) bool { return a[i].ID < a[j].ID }
