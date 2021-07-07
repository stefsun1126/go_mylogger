/* 用來打印到終端的logger */
package log

import (
	"fmt"
	"time"
)

// Logger... 日誌對象
type ConsoleLogger struct {
	// 用於限制輸出層級 大於這個值 才能輸出日誌
	OpenLevel OpenLevel
}

// NewConsoleLogger... 初始日誌對象
func NewConsoleLogger(str string) *ConsoleLogger {
	level := parseStrToOpenLevel(str)
	return &ConsoleLogger{
		OpenLevel: level,
	}
}

// 是否可在終端寫日誌訊息
func (c *ConsoleLogger) isEnable(openLevel OpenLevel) bool {
	return c.OpenLevel <= openLevel
}

func (c *ConsoleLogger) log(level OpenLevel, format string, arg ...interface{}) {
	if c.isEnable(level) {
		// main > DEBUG(假設) > log 所以有3層
		fileName, funcName, line := getInfo(3)

		// 開放使用者傳入參數可以自行格式化
		// 重組字串
		msg := fmt.Sprintf(format, arg...)
		fmt.Printf("[%v] [%v] [%v:%v:%v] %v\n", time.Now().Format("2006-01-02 15:04:05"), parseOpenLevelToStr(level), fileName, funcName, line, msg)
	}
}

func (c *ConsoleLogger) Debug(format string, arg ...interface{}) {
	c.log(DEBUG, format, arg...)
}

func (c *ConsoleLogger) Trace(format string, arg ...interface{}) {
	c.log(TRACE, format, arg...)
}

func (c *ConsoleLogger) Info(format string, arg ...interface{}) {
	c.log(INFO, format, arg...)
}

func (c *ConsoleLogger) Warning(format string, arg ...interface{}) {
	c.log(WARNING, format, arg...)
}

func (c *ConsoleLogger) Error(format string, arg ...interface{}) {
	c.log(ERROR, format, arg...)
}

func (c *ConsoleLogger) Fatal(format string, arg ...interface{}) {
	c.log(FATAL, format, arg...)
}
