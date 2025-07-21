package logger

import (
	"log"
	"os"
)

// Logger interface for dependency injection
type Logger interface {
	Info(message string)
	Error(err error)
	Fatal(err error)
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type DefaultLogger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func New() Logger {
	return &DefaultLogger{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *DefaultLogger) Info(message string) {
	l.infoLog.Println(message)
}

func (l *DefaultLogger) Error(err error) {
	l.errorLog.Println(err)
}

func (l *DefaultLogger) Fatal(err error) {
	l.errorLog.Fatal(err)
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	l.infoLog.Printf(format, args...)
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	l.errorLog.Printf(format, args...)
}