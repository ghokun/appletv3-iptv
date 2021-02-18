package logging

import (
	"io"
	"log"
	"os"
	"path"
	"time"
)

var (
	logToFile         bool
	loggingPath       string
	latestLogFilePath string
	logFile           *os.File
	Logger            *log.Logger
)

func init() {
	Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func EnableLoggingToFile(loggingPathPtr *string) {
	logToFile = true
	loggingPath = *loggingPathPtr
	t := time.Now().Format("2006-01-02")
	latestLogFilePath = path.Join(loggingPath, t+".log")
	err := os.Mkdir(loggingPath, os.ModePerm)
	logFile, err := os.OpenFile(latestLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	Logger.SetOutput(mw)
}

func CheckLogRotationAndRotate() {
	if !logToFile {
		return
	}
	t := time.Now().Format("2006-01-02")
	newLogFilePath := path.Join(loggingPath, t+".log")

	if latestLogFilePath != newLogFilePath {
		Logger.Println("Rotating log file")
		logFile.Close()
		err := os.Mkdir(loggingPath, os.ModePerm)
		logFile, err := os.OpenFile(newLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		mw := io.MultiWriter(os.Stdout, logFile)
		Logger.SetOutput(mw)
		Logger.Println("Rotated log file")
	}
}

func logInternal(prefix string, data interface{}) {
	Logger.SetPrefix(prefix)
	Logger.Println(data)
}

func Info(data interface{}) {
	logInternal("INFO: ", data)
}

func Warn(data interface{}) {
	logInternal("WARN: ", data)
}

func Fatal(data interface{}) {
	logInternal("FATAL: ", data)
	logFile.Close()
	os.Exit(1)
}
