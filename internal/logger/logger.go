package logger

import (
	"log"
	"os"
	"runtime"
)

type LogLevel int

const (
	InfoLevel LogLevel = iota
	WarningLevel
	ErrorLevel
)

type Logger struct {
	Level         LogLevel
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

var logger *Logger

func GetLogger() *Logger {
	if logger == nil {
		logger = &Logger{
			Level:         InfoLevel,
			infoLogger:    log.New(os.Stdout, "INFO: ", log.LstdFlags),
			warningLogger: log.New(os.Stdout, "WARNING: ", log.LstdFlags),
			errorLogger:   log.New(os.Stderr, "ERROR: ", log.LstdFlags),
		}
	}
	return logger
}

func (logger *Logger) SetLogLevel(level LogLevel) {
	logger.Level = level
}

func (logger *Logger) Info(msg string) {
	if logger.Level <= InfoLevel {
		logger.infoLogger.Println(msg)
	}
}

func (logger *Logger) Warning(msg string) {
	if logger.Level <= WarningLevel {
		logger.warningLogger.Println(msg)
	}
}

func (logger *Logger) Error(msg string) {
	if logger.Level <= ErrorLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			logger.errorLogger.Printf("%s:%d: %s", file, line, msg)
		} else {

			logger.errorLogger.Println(msg)
		}
	}
}
