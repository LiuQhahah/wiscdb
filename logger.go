package main

import (
	"log"
	"os"
)

type Logger interface {
	Errorf(string, ...interface{})
	Warningf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}

type defaultLog struct {
	*log.Logger
	level loggingLevel
}

type loggingLevel int

func defaultLogger(level loggingLevel) *defaultLog {
	return &defaultLog{
		Logger: log.New(os.Stderr, "wisc ", log.LstdFlags),
		level:  level,
	}
}

const (
	DEBUG loggingLevel = iota
	INFO
	WARNING
	ERROR
)

func (l *defaultLog) Errorf(f string, v ...interface{}) {

	if l != nil && l.level <= ERROR {
		log.Printf("ERROR: "+f, v...)
	}
}
func (l *defaultLog) Warningf(f string, v ...interface{}) {
	if l.level <= WARNING {
		l.Printf("WARNING: "+f, v...)
	}
}
func (l *defaultLog) Infof(f string, v ...interface{}) {
	if l.level <= INFO {
		l.Printf("INFO:"+f, v...)
	}
}

func (l *defaultLog) Debugf(f string, v ...interface{}) {
	if l.level < DEBUG {
		log.Printf("DEBUG:"+f, v...)
	}
}
