package app

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"time"
)

// PublishFiles : Publish the files into db periodically
func (p *Pipeline) PublishFiles() {
	p.Logger.Infow("Stage: PublishFiles")
	timestamp := time.Now()
	batchId := uuid.NewString()
	for _, folder := range CsvFolderLocations {
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				go p.DB.Write(RethinkDBFileTable, File{
					Timestamp: timestamp,
					BatchId:   batchId,
					Filename:  info.Name(),
					Status:    InProgress,
					Stage:     ReadyForRead,
					FilePath:  path,
				})
			}
			return nil
		})
		if err != nil {
			p.Logger.Error(err)
		}
	}
	p.Logger.Infow("Stage: PublishFiles Completed")
}

// SubscribeChangesWhenReadyForRead : Subscribe the changes when the file is ready for read.
// Read the last entry of the CSV file.
// Update the file with last entry into db.
// Update the status of the file.
func (p *Pipeline) SubscribeChangesWhenReadyForRead() {
	p.Logger.Infow("Stage: SubscribeChangesWhenReadyForRead")
	byStage, err := p.DB.ReadChangesByStage(RethinkDBFileTable, ReadyForRead)
	if err != nil {
		p.Logger.Error("Failed while reading the data by stage", err)
		return
	}

	go func() {
		var file File
		for byStage.Next(&file) {
			if len(file.ID) != 0 {
				p.Logger.Infow("getByReadyForReadStage", "file", file)
				csv := CSV{
					Logger: p.Logger,
					DB:     p.DB,
				}
				go csv.ReadLastEntryOfCsv(file)
			}
		}
	}()
}

// SubscribeChangesWhenReadyForTransform : Subscribe the changes when the file is ready for transform.
// Read the file when the stage is ReadyForTransform.
// Transform the values in TimeSeriesData format. Save into different db table
// Update the status of the file.
func (p *Pipeline) SubscribeChangesWhenReadyForTransform() {
	p.Logger.Infow("Stage: SubscribeChangesWhenReadyForTransform")
	byStage, err := p.DB.ReadChangesByStage(RethinkDBFileTable, ReadyForTransform)
	if err != nil {
		p.Logger.Error("Failed while reading the data by stage", err)
		return
	}

	go func() {
		var file File
		for byStage.Next(&file) {
			if len(file.ID) != 0 && file.LastEntry != nil {
				p.Logger.Infow("getByReadyForTransformStage", "file", file)
				transformer := Transformer{
					Logger: p.Logger,
					DB:     p.DB,
				}
				go transformer.toTimeSeries(file)
			}
		}
	}()
}

// SubscribeChangesWhenReadyForArchive : Subscribe the changes when the file is ready for archive.
// Read the file when the stage is ReadyForArchive.
// Archive the file to different folder
// Update the status of the file.
func (p *Pipeline) SubscribeChangesWhenReadyForArchive() {
	p.Logger.Infow("Stage: SubscribeChangesWhenReadyForArchive")
	byStage, err := p.DB.ReadChangesByStage(RethinkDBFileTable, ReadyForArchive)
	if err != nil {
		p.Logger.Error("Failed while reading the data by stage", err)
		return
	}

	go func() {
		var file File
		for byStage.Next(&file) {
			if len(file.ID) != 0 {
				p.Logger.Infow("getByReadyForArchiveStage", "file", file)
				err = os.Rename(file.FilePath, fmt.Sprintf(ArchivedPath, file.Filename))
				if err != nil {
					p.Logger.Errorw("Error while archiving the file", "error", err)
				}
				file.Status = Archived
				file.Stage = Completed
				go p.DB.Write(RethinkDBFileTable, file)
			}
		}
	}()
}

// Invoke : invoke the pipeline
func (p *Pipeline) Invoke() {
	p.SubscribeChangesWhenReadyForRead()
	p.SubscribeChangesWhenReadyForTransform()
	p.SubscribeChangesWhenReadyForArchive()
}
