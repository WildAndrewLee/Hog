package logger

import (
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	handler := os.Stdout
	params := log.Ldate | log.Ltime | log.Lshortfile
	Info = log.New(handler, "INFO: ", params)
	Warning = log.New(handler, "WARNING: ", params)
	Error = log.New(handler, "ERROR: ", params)
}
