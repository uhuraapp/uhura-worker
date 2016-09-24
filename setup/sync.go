package setup

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jrallison/go-workers"
	"github.com/uhuraapp/uhura-api/database"
	runner "github.com/uhuraapp/uhura-worker/sync"
)

func sync(message *workers.Msg) {
	defer reporter(message)

	id, err := message.Args().Int64()
	checkError(err)

	p := database.NewPostgresql()
	defer func(p *gorm.DB) {
		if r := recover(); r != nil {
			p.Close()
		}
	}(p)

	_, model := runner.Sync(id, p)
	p.Close()

	nextRunAt, err := runner.GetNextRun(model)

	if err != nil {
		nextRunAt = time.Now().Add(5 * time.Hour)
	}

	workers.EnqueueAt("sync", "sync", nextRunAt, id)
}
