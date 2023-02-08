package app

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

func (j *JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("jsonb type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}
