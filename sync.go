package main

import (
	"bitbucket.org/dukex/uhura-api/database"
	runner "bitbucket.org/dukex/uhura-worker/sync"
	"github.com/jrallison/go-workers"
)

func sync(message *workers.Msg) {
	defer reporter(message)

	id, err := message.Args().Int64()
	checkError(err)

	syncer(id)
	// 		nextRunAt, err := s.GetNextRun()
	// 		checkError(err)
	//
	// 		workers.EnqueueAt("sync", "sync", nextRunAt, id)
	// 		workers.Enqueue("orphan-channel", "orphanChannel", nil)
}

func syncLow(message *workers.Msg) {
	defer reporter(message)

	id, err := message.Args().Int64()
	checkError(err)

	syncer(id)
}

func syncer(id int64) {
	p := database.NewPostgresql()
	runner.Sync(id, p)
	p.Close()
}
