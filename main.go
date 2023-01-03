package main

import (
	"github.com/go-co-op/gocron"
	"time"
	"we-got-pipeline/app"
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
	s.Every(app.ScheduleInterval).Hours().Do(func() {
		p.PublishFiles()
	})
	s.StartBlocking()
}
