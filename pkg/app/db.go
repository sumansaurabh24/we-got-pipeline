package app

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
)

// InitializeDb : Initialize the database with required values
func InitializeDb(logger *zap.SugaredLogger) (*gorm.DB, error) {
	dsn := GetPostgresDSN()
	logger.Infow("Initializing database connection with connection url", "dsn", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, err
}

// Write : write the entity to the database
func (d *DBClient[E]) Write(e E) *E {
	d.Logger.Infow("Writing into table", "entity", e, "type", reflect.TypeOf(e))
	result := d.DB.Create(&e)
	if result.Error != nil {
		d.Logger.Errorw("Error while writing in database", "error", result.Error, "entity", e, "type", reflect.TypeOf(e))
		return nil
	}
	d.Logger.Infow("Write Successful", "result", result.RowsAffected, "type", reflect.TypeOf(e))
	return &e
}

// Update : Update the entity into database
func (d *DBClient[E]) Update(e E) *E {
	d.Logger.Infow("Updating the record", "entity", e, "type", reflect.TypeOf(e))
	result := d.DB.Save(&e)
	if result.Error != nil {
		d.Logger.Errorw("Error while updating database", "error", result.Error, "type", reflect.TypeOf(e))
		return nil
	}
	d.Logger.Infow("Update Successful", "result", result.RowsAffected, "type", reflect.TypeOf(e))
	return &e
}
