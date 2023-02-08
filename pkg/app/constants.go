package app

import "fmt"

const (
	PostgresHost          = "localhost"
	PostgresPort          = 5432
	PostgresDBName        = "wegot"
	PostgresUserName      = "postgres"
	PostgresPassword      = ""
	PostgresConnectionDSN = "host=%s user=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Calcutta"
	ArchivedPath          = "/Users/suman.saurabh/Archived/%s"
	EnableArchiving       = false
	ScheduleInterval      = 24
	RegExpPattern         = "(A|B|C|D)[0-9]+(\\s)+_([A-Z][0-9]|[A-Z]+)"
	FilenamePrefix        = "VENAQUA_"
)

var CsvFolderLocations = [...]string{"/Users/suman.saurabh/Project/data/a", "/Users/suman.saurabh/Project/data/b"}

const (
	InProgress Status = "IN_PROGRESS"
	Success    Status = "SUCCESS"
	Archived   Status = "ARCHIVED"
)

const (
	ReadyForRead      Stage = "READY_FOR_READ"
	ReadyForTransform Stage = "READY_FOR_TRANSFORM"
	ReadyForArchive   Stage = "READY_FOR_ARCHIVE"
	Completed         Stage = "COMPLETED"
)

func GetPostgresDSN() string {
	return fmt.Sprintf(
		PostgresConnectionDSN, PostgresHost, PostgresUserName, PostgresDBName, PostgresPort,
	)
}
