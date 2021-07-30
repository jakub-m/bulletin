package log

import "log"

func Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
