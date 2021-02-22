package logging

import (
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/ghokun/appletv3-iptv/internal/config"
)

var (
	latestLogFilePath string
	logFile           *os.File
	Logger            *log.Logger
)

func init() {
	Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func EnableLoggingToFile() {
	t := time.Now().Format("2006-01-02")
	latestLogFilePath = path.Join(config.Current.LoggingPath, t+".log")
	err := os.Mkdir(config.Current.LoggingPath, os.ModePerm)
	logFile, err := os.OpenFile(latestLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	Logger.SetOutput(mw)
}

func CheckLogRotationAndRotate() {
	if !config.Current.LogToFile {
		return
	}
	t := time.Now().Format("2006-01-02")
	newLogFilePath := path.Join(config.Current.LoggingPath, t+".log")

	if latestLogFilePath != newLogFilePath {
		Logger.Println("Rotating log file")
		logFile.Close()
		err := os.Mkdir(config.Current.LoggingPath, os.ModePerm)
		logFile, err := os.OpenFile(newLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
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
