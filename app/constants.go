package app

const (
	MongoDBHost            = "mongodb://localhost:%d/"
	MongoDBPort            = 27017
	MongoDBName            = "we-got"
	MongoDBFileTable       = "file"
	MongoDBTimeseriesTable = "timeseries"
	ArchivedPath           = "/Users/suman.saurabh/Archived/%s"
	ScheduleInterval       = 24
	RegExpPattern          = "(A|B|C|D)[0-9]+(\\s)+_([A-Z][0-9]|[A-Z]+)"
	FilenamePrefix         = "VENAQUA_"
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
