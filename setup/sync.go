package setup

import (
	"time"

	"github.com/uhuraapp/uhura-api/database"
	runner "github.com/uhuraapp/uhura-worker/sync"
	"github.com/jinzhu/gorm"
	"github.com/jrallison/go-workers"
)

func sync(message *workers.Msg) {
	defer reporter(message)

	id, err := message.Args().Int64()
	checkError(err)

	syncer(id, true)
}

func syncLow(message *workers.Msg) {
	defer reporter(message)

	id, err := message.Args().Int64()
	checkError(err)

	syncer(id, false)

	workers.EnqueueAt("sync-low", "syncLow", time.Now().Add(1*time.Hour), id)
}

func syncer(id int64, scheduleNext bool) {
	p := database.NewPostgresql()
	defer func(p gorm.DB) {
		if r := recover(); r != nil {
			p.Close()
		}
	}(p)

	_, model := runner.Sync(id, p)
	p.Close()

	workers.Enqueue("duplicate-episodes", "duplicateEpisodes", nil)
	if scheduleNext {
		nextRunAt, err := runner.GetNextRun(model)
		checkError(err)
		workers.EnqueueAt("sync", "sync", nextRunAt, id)
	}
}
