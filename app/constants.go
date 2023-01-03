package app

const (
	RethinkDBHost            = "localhost:%d"
	RethinkDBPort            = 28015
	RethinkDBName            = "we-got"
	RethinkDBFileTable       = "file"
	RethinkDBTimeseriesTable = "timeseries"
	RethinkDBChangeNewValue  = "new_val"
	ArchivedPath             = "/Users/suman.saurabh/Archived/%s"
	ScheduleInterval         = 24
	RegExpPattern            = "(A|B|C|D)[0-9]+(\\s)+_([A-Z][0-9]|[A-Z]+)"
	FilenamePrefix           = "VENAQUA_"
	MetadataTotalKey         = "total"
)

var CsvFolderLocations = [...]string{"/Users/suman.saurabh/Project/data/a", "/Users/suman.saurabh/Project/data/b"}

const (
	InProgress Status = "IN_PROGRESS"
	Success    Status = "SUCCESS"
	Failed     Status = "FAILED"
	Archived   Status = "ARCHIVED"
)

const (
	ReadyForRead      Stage = "ReadyForRead"
	ReadyForTransform Stage = "ReadyForTransform"
	ReadyForArchive   Stage = "ReadyForArchive"
	Completed         Stage = "Completed"
)
