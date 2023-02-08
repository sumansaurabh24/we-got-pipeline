package app

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Pipeline struct {
	Logger                   *zap.SugaredLogger
	PubSub                   *PubSub
	FileRepository           *DBClient[File]
	TimeSeriesDataRepository *DBClient[TimeSeriesData]
}

type Entity interface {
	File | TimeSeriesData
}

type DBClient[E Entity] struct {
	Logger *zap.SugaredLogger
	DB     *gorm.DB
}

type Status string

type Stage string

type File struct {
	ID        uint      `gorm:"column:id"`
	BatchId   string    `gorm:"column:batch_id""`
	Filename  string    `gorm:"column:filename"`
	FilePath  string    `gorm:"column:file_path"`
	Status    Status    `gorm:"column:status"`
	Stage     Stage     `gorm:"column:stage"`
	Timestamp time.Time `gorm:"column:timestamp;autoCreateTime"`
	LastEntry JSONB     `gorm:"column:last_entry;type:jsonb"`
}

type TimeSeriesData struct {
	FilePath         string    `gorm:"column:file_path"`
	Flat             string    `gorm:"column:flat_no"`
	Metadata         JSONB     `gorm:"column:metadata;type:jsonb"`
	TotalConsumption int64     `gorm:"column:total_consumption"`
	Timestamp        time.Time `gorm:"column:timestamp;"`
}

type CSV struct {
	Logger *zap.SugaredLogger
}

type Transformer struct {
	Logger                   *zap.SugaredLogger
	PubSub                   *PubSub
	FileRepository           *DBClient[File]
	TimeSeriesDataRepository *DBClient[TimeSeriesData]
}

type JSONB map[string]interface{}

type LumberjackSink struct {
	*lumberjack.Logger
}

func (l LumberjackSink) Sync() error {
	return nil
}
