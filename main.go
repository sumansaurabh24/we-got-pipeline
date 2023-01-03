package main

import (
	"data-pipeline/app"
	"github.com/go-co-op/gocron"
	"time"
)

func main() {
	logger, err := app.InitializeLogger()
	if err != nil {
		panic(err)
	}
	db, err := app.InitializeDb(logger)
	if err != nil {
		logger.Panic(err)
	}

	p := app.Pipeline{
		Logger: logger,
		DB:     db,
	}
	p.Invoke()

	s := gocron.NewScheduler(time.UTC)
	s.Every(10).Hours().Do(func() {
		p.PublishFiles()
	})
	s.StartBlocking()
}
