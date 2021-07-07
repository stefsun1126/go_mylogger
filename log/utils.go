package log

import (
	"path"
	"runtime"
	"strings"
)

// 放共用
// 日誌層級輸出開關用
type OpenLevel uint16

// 日誌層級輸出常數
const (
	UNKNOW OpenLevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

// log interface
type MyLogger interface {
	Debug(format string, arg ...interface{})
	Info(format string, arg ...interface{})
	Warning(format string, arg ...interface{})
	Error(format string, arg ...interface{})
	Fatal(format string, arg ...interface{})
}

// OpenLevel => str to OpenLevel
func parseStrToOpenLevel(str string) OpenLevel {
	var level OpenLevel
	switch strings.ToLower(str) {
	case "debug":
		level = DEBUG
	case "trace":
		level = TRACE
	case "info":
		level = INFO
	case "warning":
		level = WARNING
	case "error":
		level = ERROR
	case "fatal":
		level = FATAL
	default:
		level = UNKNOW
	}
	return level
}

// // OpenLevel => OpenLevel to str
func parseOpenLevelToStr(openLevel OpenLevel) string {
	var str string
	switch openLevel {
	case DEBUG:
		str = "DEBUG"
	case TRACE:
		str = "TRACE"
	case INFO:
		str = "INFO"
	case WARNING:
		str = "WARNING"
	case ERROR:
		str = "ERROR"
	case FATAL:
		str = "FATAL"
	default:
		str = "UNKNOW"
	}
	return str
}

// 獲取執行中的文件名、函數、行號
func getInfo(skip int) (fileName, funcName string, line int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		panic("runtime.Caller fail!")
	}
	// 獲取執行中的函數名
	funcName = strings.Split(runtime.FuncForPC(pc).Name(), ".")[1]
	fileName = path.Base(file)
	return
}
