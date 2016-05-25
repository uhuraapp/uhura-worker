package setup

import (
	"github.com/uhuraapp/uhura-api/database"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/jrallison/go-workers"
)

func deleteEpisode(message *workers.Msg) {
	id, err := message.Args().Int64()
	checkError(err)

	p := database.NewPostgresql()
	p.Table(models.Episode{}.TableName()).Where("id = ?", id).Delete(models.Episode{})
	p.Close()
}
