package main

import (
	"log"
	"os"
)

type loggers struct {
	info *log.Logger
	err  *log.Logger
}

var logging *loggers

func InitLoggers() {
	info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	err := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	logging = &loggers{info, err}
}

func LogInfo(msg string) {
	logging.info.Println(msg)
}

func LogError(msg string) {
	logging.err.Println(msg)
}
