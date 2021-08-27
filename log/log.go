package log

import (
	golog "log"
)

var verboseMode = false

func SetVerbose(verbose bool) {
	verboseMode = verbose
}

func Infof(format string, args ...interface{}) {
	golog.Printf("I "+format, args...)
}

func Debugf(format string, args ...interface{}) {
	if verboseMode {
		golog.Printf("D "+format, args...)
	}
}
