package app

import (
	"context"
	"fmt"
	"github.com/naamancurtis/mongo-go-struct-to-bson/mapper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// InitializeDb : Initialize the database with required values
func InitializeDb(logger *zap.SugaredLogger) (*DB, error) {
	logger.Infow("Initializing database connection")
	clientOptions := options.Client().ApplyURI(fmt.Sprintf(MongoDBHost, MongoDBPort))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	db := &DB{
		Logger: logger,
		Client: client,
	}
	return db, err
}

func (d *DB) Collection(name string) *mongo.Collection {
	d.Logger.Infow("Collection", "name", name)
	return d.Client.Database(MongoDBName).Collection(name)
}

// Write : Write into database tables and return inserted id
func (d *DB) Write(table string, data interface{}) *primitive.ObjectID {
	d.Logger.Infow("Writing into table", "table_name", table, "data", data)
	result, err := d.Collection(table).InsertOne(context.TODO(), data)
	if err != nil {
		d.Logger.Errorw("Error while writing in database", "error", err, "table", table)
		return nil
	}
	d.Logger.Infow("Write Successful", "result", result.InsertedID)
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return &oid
	}
	return nil
}

// UpdateById : Update the document data by its id and return id
func (d *DB) UpdateById(table string, file File) *primitive.ObjectID {
	d.Logger.Infow("Updating the record", "table_name", table, "data", file)
	filter := bson.D{{"_id", file.ID}}
	update := bson.D{{"$set", mapper.ConvertStructToBSONMap(file, &mapper.MappingOpts{RemoveID: true})}}
	result, err := d.Collection(table).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		d.Logger.Errorw("Error while updating database", "error", err, "table", table)
		return nil
	}
	d.Logger.Infow("Update Successful", "result", result.ModifiedCount)
	if result.ModifiedCount > 0 {
		return file.ID
	}

	return nil
}
