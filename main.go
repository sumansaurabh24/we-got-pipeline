package main

import (
	"context"
	"github.com/go-co-op/gocron"
	"time"
	"we-got-pipeline/app"
)

func main() {
	//###################### Initialization ###############################
	logger, err := app.InitializeLogger()
	if err != nil {
		panic(err)
	}
	db, err := app.InitializeDb(logger)
	if err != nil {
		logger.Panic(err)
	}
	defer db.Client.Disconnect(context.TODO())

	pubSub := app.InitializeNats(logger)
	//###################### Initialization ###############################

	p := app.Pipeline{
		Logger: logger,
		DB:     db,
		PubSub: pubSub,
	}
	p.Invoke()

	s := gocron.NewScheduler(time.UTC)
	s.Every(app.ScheduleInterval).Hours().Do(func() {
		p.PublishFiles()
	})
	s.StartBlocking()
}
