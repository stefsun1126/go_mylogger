package main

import (
	"mylogger/log"
)

var (
	logger      log.MyLogger
	fileName    string = "file_logger"
	maxFileSize int    = 10 * 1024 * 1024
)

func main() {
	// logger = log.NewConsoleLogger("debug")
	logger = log.NewFileLog("debug", fileName, "./", maxFileSize)

	for {
		logger.Debug("這是Debug層級的log id=%v test=%v", "1234", "28825252")
		logger.Info("這是Info層級的log")
		logger.Warning("這是Warning層級的log")
		logger.Error("這是Error層級的log")
		logger.Fatal("這是Fatal層級的log")
		// time.Sleep(3 * time.Second)
	}

}
