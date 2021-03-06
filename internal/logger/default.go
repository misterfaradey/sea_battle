package logger

import (
	"log"
	"os"
)

var defaultLogger = log.New(os.Stdout, "default_logger", log.Ldate|log.Ltime|log.Lshortfile)

func DefaultLogger() *log.Logger {
	return defaultLogger
}
