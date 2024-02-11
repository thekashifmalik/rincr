package internal

var NAME = "rincr"
var BACKUPS_DIR = "." + NAME
var BACKUPS_DIR_PATH = "/" + BACKUPS_DIR
var LAST_FILE_PATH = BACKUPS_DIR_PATH + "/last"

// TIME_FORMAT is the format used to in the backup folder names. The format needs to be file-name-safe. We try to match
// the universally parse-able ISO format as close as possible.
var TIME_FORMAT = "2006-01-02T15-04-05"
