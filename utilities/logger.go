/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package utilities

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	stdoutOnce         sync.Once
	logsingleton       *log.Logger
	primaryLogFile     string
	stdoutLogSingleton *log.Logger
)

// CreateLogger create a new logger with given configuration
// isDebug = true, will set the log level to be Debug
// maxSize is the maximum size of the file (in MB) before it is rotated
// maxBackups is the number of backups to keep
// maxAge is the number of days after which log file will be deleted
// fileName is the name of log file. If fileName is empty, it will create logger for stdout.
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

// GetLogger returns the (already created) logger
func GetLogger() *log.Logger {
	if logsingleton != nil {
		return logsingleton
	}
	return stdoutLogSingleton
}
