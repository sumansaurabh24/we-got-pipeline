package main

import (
	"github.com/go-co-op/gocron"
	"time"
	"we-got-pipeline/pkg/app"
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
	//db.AutoMigrate(&app.File{}, &app.TimeSeriesData{})
	pubSub := app.InitializeNats(logger)
	//###################### Initialization ###############################

	p := app.Pipeline{
		Logger: logger,
		PubSub: pubSub,
		FileRepository: &app.DBClient[app.File]{
			Logger: logger,
			DB:     db,
		},
		TimeSeriesDataRepository: &app.DBClient[app.TimeSeriesData]{
			Logger: logger,
			DB:     db,
		},
	}
	p.Invoke()

	s := gocron.NewScheduler(time.UTC)
	s.Every(app.ScheduleInterval).Hours().Do(func() {
		p.PublishFiles()
	})
	s.StartBlocking()
}
