package app

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"os"
	"path/filepath"
	"time"
)

// insertFileAndPublish : insert into database and publish to message queue
func (p *Pipeline) insertFileAndPublish(file File) {
	insertedID := p.DB.Write(MongoDBFileTable, file)
	if insertedID != nil {
		file.ID = insertedID
		go p.PubSub.Publish(ReadyForRead, file)
	}
}

// PublishFiles : Publish the files into db periodically
func (p *Pipeline) PublishFiles() {
	p.Logger.Infow("Stage: PublishFiles")
	timestamp := time.Now()
	batchId := uuid.NewString()
	for _, folder := range CsvFolderLocations {
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				file := File{
					Timestamp: timestamp,
					BatchId:   batchId,
					Filename:  info.Name(),
					Status:    InProgress,
					Stage:     ReadyForRead,
					FilePath:  path,
				}
				go p.insertFileAndPublish(file)
			}
			return nil
		})
		if err != nil {
			p.Logger.Error(err)
		}
	}
	p.Logger.Infow("Stage: PublishFiles Completed")
}

// readyForReadSubscriber : Subscribe the changes when the file is ready for read.
// Read the last entry of the CSV file.
// Update the file with last entry into db.
// Update the status of the file.
func (p *Pipeline) readyForReadSubscriber(msg *nats.Msg) {
	p.Logger.Infow("Stage: readyForReadSubscriber")

	go func() {
		var file File
		if err := json.Unmarshal(msg.Data, &file); err != nil {
			p.Logger.Error("Error while decoding", err)
			return
		}
		p.Logger.Infow("readyForReadSubscriber received message", "file", file)
		csv := CSV{
			Logger: p.Logger,
		}
		lastEntry := csv.ReadLastEntryOfCsv(file)

		//change few values before updating
		file.LastEntry = lastEntry
		file.Stage = ReadyForTransform
		p.Logger.Infow("Last entries map created for file", "id", file.ID)

		go func() {
			upsertedID := p.DB.UpdateById(MongoDBFileTable, file)
			if upsertedID != nil {
				go p.PubSub.Publish(ReadyForTransform, file)
			}
		}()
	}()
}

// readyForTransformSubscriber : Subscribe the changes when the file is ready for transform.
// Read the file when the stage is ReadyForTransform.
// Transform the values in TimeSeriesData format. Save into different db table
// Update the status of the file.
func (p *Pipeline) readyForTransformSubscriber(msg *nats.Msg) {
	p.Logger.Infow("Stage: readyForTransformSubscriber")

	go func() {
		var file File
		if err := json.Unmarshal(msg.Data, &file); err != nil {
			p.Logger.Error("Error while decoding", err)
			return
		}
		p.Logger.Infow("readyForTransformSubscriber received message", "file", file)
		transformer := Transformer{
			Logger: p.Logger,
			DB:     p.DB,
			PubSub: p.PubSub,
		}
		go transformer.toTimeSeries(file)
	}()
}

// readyForArchiveSubscriber : Subscribe the changes when the file is ready for archive.
// Read the file when the stage is ReadyForArchive.
// Archive the file to different folder
// Update the status of the file.
func (p *Pipeline) readyForArchiveSubscriber(msg *nats.Msg) {
	p.Logger.Infow("Stage: readyForArchiveSubscriber")

	go func() {
		var file File
		if err := json.Unmarshal(msg.Data, &file); err != nil {
			p.Logger.Error("Error while decoding", err)
			return
		}
		p.Logger.Infow("readyForArchiveSubscriber received message", "file", file)
		if err := os.Rename(file.FilePath, fmt.Sprintf(ArchivedPath, file.Filename)); err != nil {
			p.Logger.Errorw("Error while archiving the file", "error", err)
			return
		}
		file.Status = Archived
		file.Stage = Completed
		go p.DB.UpdateById(MongoDBFileTable, file)
	}()
}

// Invoke : invoke the pipeline
func (p *Pipeline) Invoke() {
	p.PubSub.RegisterSubscribers(map[string]nats.MsgHandler{
		string(ReadyForRead):      p.readyForReadSubscriber,
		string(ReadyForTransform): p.readyForTransformSubscriber,
		string(ReadyForArchive):   p.readyForArchiveSubscriber,
	})
}
