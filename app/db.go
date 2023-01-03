package app

import (
	"fmt"
	"go.uber.org/zap"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// InitializeDb : Initialize the database with required values
func InitializeDb(logger *zap.SugaredLogger) (*DB, error) {
	logger.Infow("Initializing database connection")
	session, err := r.Connect(r.ConnectOpts{
		Address:  fmt.Sprintf(RethinkDBHost, RethinkDBPort),
		Database: RethinkDBName,
	})
	db := &DB{
		Logger:  logger,
		Session: session,
	}
	return db, err
}

// Write : Write into database tables
func (d *DB) Write(table string, data interface{}) {
	d.Logger.Infow("Writing into table", "table_name", table, "data", data)
	_, err := r.Table(table).Insert(data).Run(d.Session)
	if err != nil {
		d.Logger.Error("Error while writing in database", err)
	}
}

// ReadChangesByStage : Read the streaming changes from database tables
func (d *DB) ReadChangesByStage(table string, stage Stage) (*r.Cursor, error) {
	d.Logger.Infow("Reading from ", "table_name", table, "stage_name", stage)
	cursor, err := r.Table(table).Filter(map[string]Stage{
		"stage": stage,
	}).Changes().Field(RethinkDBChangeNewValue).Run(d.Session)
	return cursor, err
}

// Update : Update the table data
func (d *DB) Update(table string, data interface{}) {
	d.Logger.Infow("Updating the record", "table_name", table, "data", data)
	_, err := r.Table(table).Update(data).Run(d.Session)
	if err != nil {
		d.Logger.Error("Error while updating database", err)
	}
}
