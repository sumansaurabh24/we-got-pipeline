package app

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// read : Read csv file with all data
func (c *CSV) read(filepath string) [][]string {
	file, err := os.Open(filepath)
	if err != nil {
		c.Logger.Errorw("Error while opening the file", "filename", filepath, err)
		return nil
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()

	if err != nil {
		c.Logger.Errorw("Error while reading all in csv reader", err)
		return nil
	}

	log.Println(fmt.Sprintf("Length of array %d", len(data)))
	return data
}

// ReadLastEntryOfCsv : Read all csv content and extract last entry/tuple/row of the csv
func (c *CSV) ReadLastEntryOfCsv(file File) {
	data := c.read(file.FilePath)
	if data == nil {
		c.Logger.Warnw("Data is nil after reading csv")
		return
	}

	rowLength := len(data)
	if rowLength < 2 {
		c.Logger.Warnw("Row length is less than 2 after reading csv", "row_length", rowLength)
		return
	}

	columnLength := len(data[0])
	if columnLength < 2 {
		c.Logger.Warnw("Column length is less than 2 after reading csv: %d", "column_length", columnLength)
		return
	}

	lastRowIndex := rowLength - 1
	lastEntry := make(map[string]string)
	for i := 0; i < columnLength; i++ {
		headerValue := data[0][i]
		if len(headerValue) == 0 {
			headerValue = fmt.Sprintf("__EMPTY_%d", i)
		}
		lastEntry[headerValue] = data[lastRowIndex][i]
	}

	//change few values before updating
	file.LastEntry = lastEntry
	file.Stage = ReadyForTransform

	go c.DB.Update(RethinkDBFileTable, file)
	c.Logger.Infow("Last entries map created for file", "id", file.ID)
}
