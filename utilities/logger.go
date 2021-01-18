package utilities

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger = *log.Logger

var (
	stdoutOnce         sync.Once
	logsingleton       *log.Logger
	primaryLogFile     string
	stdoutLogSingleton *log.Logger
)

func CreateLogger(isDebug bool, maxSize int, maxBackups int, maxAge int, fileName ...string) *log.Logger {

	if len(fileName) == 0 && logsingleton == nil {
		if stdoutLogSingleton == nil {
			stdoutOnce.Do(func() {
				stdoutLogSingleton = log.New()
				stdoutLogSingleton.SetOutput(os.Stdout)
				if isDebug {
					stdoutLogSingleton.SetLevel(log.DebugLevel)
				} else {
					stdoutLogSingleton.SetLevel(log.InfoLevel)
				}
			})
		}
		return stdoutLogSingleton
	} else if logsingleton == nil {
		logsingleton = createLogger(isDebug, maxSize, maxBackups, maxAge, fileName[0])
		primaryLogFile = fileName[0]
	}
	return logsingleton
}

func createLogger(isDebug bool, maxSize int, maxBackups int, maxAge int, logFile string) *log.Logger {
	var newInstance *log.Logger

	newInstance = log.New()
	rotateLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxSize, // megabytes
		MaxBackups: maxBackups,
		MaxAge:     maxAge, //days
		Compress:   true,   // disabled by default
	}
	newInstance.SetOutput(rotateLogger)
	if isDebug {
		newInstance.SetLevel(log.DebugLevel)
	} else {
		newInstance.SetLevel(log.InfoLevel)
	}
	newInstance.Info("Logger created")
	return newInstance
}

func GetLogger() *log.Logger {
	if logsingleton != nil {
		return logsingleton
	} else {
		return stdoutLogSingleton
	}
}
