package app

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// toTimeSeries : Transform the linear json to time series structure
func (t *Transformer) toTimeSeries(file File) {
	var processedMap = make(map[string]TimeSeriesData)
	timestamp, err := t.getTimestampFromFilename(file.Filename)
	if err != nil {
		t.Logger.Errorw("Error while parsing time from filename", "error", err, "file_id", file.ID)
		return
	}
	for index, value := range file.LastEntry {
		matched, err := regexp.Match(RegExpPattern, []byte(index))
		if err != nil {
			t.Logger.Errorw("There was some error in pattern match, hence continuing", "Error", err)
			continue
		}
		if !matched {
			t.Logger.Warnw("Was not able to match the pattern, hence continuing", "index", index)
			continue
		}
		index = strings.ReplaceAll(index, " ", "")
		splitIndex := strings.Split(index, "_")
		intValue, _ := strconv.Atoi(value.(string))
		if v, found := processedMap[splitIndex[0]]; found {
			v.Metadata[splitIndex[1]] = intValue
		} else {
			processedMap[splitIndex[0]] = TimeSeriesData{
				FilePath:  file.FilePath,
				Flat:      splitIndex[0],
				Metadata:  map[string]interface{}{splitIndex[1]: intValue},
				Timestamp: timestamp,
			}
		}
	}

	go func() {
		for _, ts := range processedMap {
			total := 0
			for _, value := range ts.Metadata {
				total += value.(int)
			}
			ts.TotalConsumption = int64(total)
			go t.TimeSeriesDataRepository.Write(ts)
		}

		file.Stage = ReadyForArchive
		file.Status = Success
		updatedID := t.FileRepository.Update(file)
		if updatedID != nil {
			go t.PubSub.Publish(ReadyForArchive, file)
		}
	}()

}

// getTimestampFromFilename: extract time from filename and then convert it into time object
func (t *Transformer) getTimestampFromFilename(filename string) (time.Time, error) {
	indexOfDot := strings.LastIndex(filename, ".")
	exactFilenameWithoutExtension := filename[0:indexOfDot]
	dateFromFilename := strings.ReplaceAll(exactFilenameWithoutExtension, FilenamePrefix, "")
	t.Logger.Infow("extracted date from filename", "date_from_filename", dateFromFilename)
	spaceSeperatedDate := strings.ReplaceAll(dateFromFilename, "-", " ")
	theTime, err := time.Parse("02 Jan 06", spaceSeperatedDate)
	return theTime, err
}
