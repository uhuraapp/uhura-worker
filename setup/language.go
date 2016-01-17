package setup

import (
	"bitbucket.org/dukex/uhura-api/database"
	"bitbucket.org/dukex/uhura-api/parser"
	"github.com/jinzhu/gorm"

	runner "bitbucket.org/dukex/uhura-worker/sync"
	"github.com/jrallison/go-workers"
)

func language(message *workers.Msg) {
	defer reporter(message)

	id, err := message.Args().Int64()
	checkError(err)

	p := database.NewPostgresql()
	defer func(p gorm.DB) {
		if r := recover(); r != nil {
			p.Close()
		}
	}(p)

	channel := runner.GetModel(id, p)
	channel.Language = parser.NormalizeLanguage(channel.Language)
	p.Table("channels").Save(channel)

	p.Close()
}
