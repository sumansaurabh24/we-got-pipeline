package app

import (
	"github.com/natefinch/lumberjack"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type DB struct {
	Logger *zap.SugaredLogger
	Client *mongo.Client
}

type DBInterface interface {
	Collection(name string) *mongo.Collection
	Write(table string, data interface{}) *primitive.ObjectID
	UpdateById(table string, file File) *primitive.ObjectID
}

type Pipeline struct {
	Logger *zap.SugaredLogger
	DB     DBInterface
	PubSub *PubSub
}

type Status string

type Stage string

type File struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty"`
	BatchId   string              `bson:"batch_id"`
	Filename  string              `bson:"filename"`
	FilePath  string              `bson:"filepath"`
	Status    Status              `bson:"status"`
	Stage     Stage               `bson:"stage"`
	Timestamp time.Time           `bson:"timestamp"`
	LastEntry map[string]string   `bson:"last_entry"`
}

type CSV struct {
	Logger *zap.SugaredLogger
}

type Transformer struct {
	Logger *zap.SugaredLogger
	DB     DBInterface
	PubSub *PubSub
}

type TimeSeriesData struct {
	ID               string         `bson:"id,omitempty"`
	FilePath         string         `bson:"file_path"`
	Flat             string         `bson:"flat"`
	Metadata         map[string]int `bson:"metadata"`
	TotalConsumption int64          `bson:"total_consumption"`
	Timestamp        time.Time      `bson:"timestamp"`
}

type LumberjackSink struct {
	*lumberjack.Logger
}

func (l LumberjackSink) Sync() error {
	return nil
}
