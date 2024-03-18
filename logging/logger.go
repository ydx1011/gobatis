package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

var Log LogFunc = DefaultLogf

var gLogTag map[int]string = map[int]string{}
var gLogLevel = INFO

type LogFunc func(level int, format string, args ...interface{})

func DefaultLogf(level int, format string, args ...interface{}) {
	if gLogLevel > level {
		return
	}

	logInfo := fmt.Sprintf(format, args...)
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	log.Printf("%s %s:%d %s", gLogTag[level], shortFile(file), line, logInfo)
	if level >= FATAL {
		os.Exit(-1)
	}
}

func shortFile(file string) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return short
}

func Info(format string, args ...interface{}) {
	Log(INFO, format, args...)
}
func Warn(format string, args ...interface{}) {
	Log(WARN, format, args...)
}
func Debug(format string, args ...interface{}) {
	Log(DEBUG, format, args...)
}
