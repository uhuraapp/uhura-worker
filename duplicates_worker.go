package main

import (
	"time"

	"bitbucket.org/dukex/uhura-api/database"
	"github.com/jrallison/go-workers"

	duplicates "bitbucket.org/dukex/uhura-worker/duplicates"
)

func duplicateEpisodes(message *workers.Msg) {
	p := database.NewPostgresql()
	del := make(chan int64)
	cl := make(chan bool)

	go func() {
		for {
			select {
			case id := <-del:
				workers.Enqueue("delete-episode", "deleteEpisode", id)
			case <-cl:
				return
			}
		}
	}()

	duplicates.Episodes(p, del, cl)

	p.Close()

	workers.EnqueueAt("duplicate-episodes", "duplicateEpisodes", time.Now().Add(time.Hour*1), nil)
}
