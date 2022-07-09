package log

import (
	"io"
	golog "log"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelSilent
)

var logLevel = LevelInfo

func SetLogLevel(lev LogLevel) {
	logLevel = lev
	if lev == LevelSilent {
		golog.SetOutput(io.Discard)
	}
}

func Infof(format string, args ...interface{}) {
	if logLevel <= LevelInfo {
		golog.Printf("I "+format, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if logLevel <= LevelDebug {
		golog.Printf("D "+format, args...)
	}
}
