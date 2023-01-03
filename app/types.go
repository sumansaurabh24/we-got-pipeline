package app

import (
	"go.uber.org/zap"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
	"time"
)

type DB struct {
	Logger  *zap.SugaredLogger
	Session *rethinkdb.Session
}

type DBInterface interface {
	Write(table string, data interface{})
	ReadChangesByStage(table string, stage Stage) (*rethinkdb.Cursor, error)
	Update(table string, data interface{})
}

type Pipeline struct {
	Logger *zap.SugaredLogger
	DB     DBInterface
}

type Status string

type Stage string

type File struct {
	ID        string            `rethinkdb:"id,omitempty"`
	BatchId   string            `rethinkdb:"batch_id"`
	Filename  string            `rethinkdb:"filename"`
	FilePath  string            `rethinkdb:"filepath"`
	Status    Status            `rethinkdb:"status"`
	Stage     Stage             `rethinkdb:"stage"`
	Timestamp time.Time         `rethinkdb:"timestamp"`
	LastEntry map[string]string `rethinkdb:"last_entry"`
}

type CSV struct {
	Logger *zap.SugaredLogger
	DB     DBInterface
}

type Transformer struct {
	Logger *zap.SugaredLogger
	DB     DBInterface
}

type TimeSeriesData struct {
	ID        string         `rethinkdb:"id,omitempty"`
	FileId    string         `rethinkdb:"file_id"`
	Flat      string         `rethinkdb:"flat"`
	Metadata  map[string]int `rethinkdb:"metadata"`
	Timestamp time.Time      `rethinkdb:"timestamp"`
}
