/* 用來打印到文件的logger */
package log

import (
	"fmt"
	"os"
	"path"
	"time"
)

var maxChanSize int = 50000

type MsgType struct {
	Time     string
	Level    OpenLevel
	FileName string
	FuncName string
	Line     int
	Msg      string
}

// Logger ... 日誌對象
type FileLogger struct {
	// 用於限制輸出層級 大於這個值 才能輸出日誌
	OpenLevel OpenLevel
	// 文件名
	FileName string
	// 文件位置
	FilePath string
	// 文件obj
	FileObj *os.File
	// 錯誤級別文件obj
	ErrorFileObj *os.File
	// 文件最大大小
	MaxFileSize int
	// 訊息管道
	MsgChan chan *MsgType
}

// 初始日誌對象
func NewFileLog(openLevel, fileName, filePath string, maxFileSize int) *FileLogger {
	level := parseStrToOpenLevel(openLevel)
	logger := &FileLogger{
		OpenLevel:   level,
		FileName:    fileName,
		FilePath:    filePath,
		MaxFileSize: maxFileSize,
		MsgChan:     make(chan *MsgType, maxChanSize),
	}
	// 初始日誌文件
	err := logger.initFile()
	if err != nil {
		panic(err)
	}
	return logger
}

// initFile ... 初始日誌文件
func (f *FileLogger) initFile() error {
	file, err := os.OpenFile(path.Join(f.FilePath, f.FileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("os.OpenFile file fail , err: %v ", err)
	}
	errfile, err := os.OpenFile(path.Join(f.FilePath, f.FileName+"_error"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("os.OpenFile errfile fail , err: %v ", err)
	}
	f.FileObj = file
	f.ErrorFileObj = errfile

	// 開goroutine在背後等待寫文件
	go f.writeFileOnBackground()

	return nil
}

// isEnable ... 是否可在終端寫日誌訊息
func (f *FileLogger) isEnable(openLevel OpenLevel) bool {
	return f.OpenLevel <= openLevel
}

// checkFileSizeValid ... 檢查日誌大小
func (f *FileLogger) checkFileSizeValid(isErrorFile bool) (bool, error) {
	// 看文件是哪一種
	file := f.FileObj
	if isErrorFile {
		file = f.ErrorFileObj
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return false, fmt.Errorf("file.Stat() error :%v", err)
	}

	return int(fileInfo.Size()) <= f.MaxFileSize, nil
}

// divideFile ... 分割文件
func (f *FileLogger) divideFile(isErrorFile bool) error {
	file := f.FileObj
	if isErrorFile {
		file = f.ErrorFileObj
	}
	// 把舊文件改名
	bakPath := path.Join(f.FilePath, file.Name()+time.Now().Format("20060102150405000"))
	// 新增新文件並與舊文件同名
	originPath := path.Join(f.FilePath, file.Name())
	// 關閉當前文件
	file.Close()
	// 將舊文件改名
	if err := os.Rename(originPath, bakPath); err != nil {
		return fmt.Errorf("os.Rename() error :%v", err)
	}

	// 新開文件
	if newFile, err := os.OpenFile(originPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err != nil {
		return fmt.Errorf("os.OpenFile error :%v", err)
	} else {
		// 修改實例
		if isErrorFile {
			f.ErrorFileObj = newFile
		} else {
			f.FileObj = newFile
		}
	}

	return nil
}

// log ... 記錄到訊息管道
func (f *FileLogger) log(level OpenLevel, format string, arg ...interface{}) {
	if f.isEnable(level) {
		// main > DEBUG(假設) > log 所以有3層
		fileName, funcName, line := getInfo(3)

		// 開放使用者傳入參數可以自行格式化
		msg := fmt.Sprintf(format, arg...)

		// 將訊息塞到管道中
		msgInstance := &MsgType{
			Time:     time.Now().Format("2006-01-02 15:04:05"),
			Level:    level,
			FileName: fileName,
			FuncName: funcName,
			Line:     line,
			Msg:      msg,
		}
		select {
		case f.MsgChan <- msgInstance:
		default: // 數據丟失
			fmt.Printf("msg missing:%v\n", *msgInstance)
		}
	}
}

// writeFileOnBackground ... 從管道取出寫文件
func (f *FileLogger) writeFileOnBackground() {
	for {
		select {
		case msg := <-f.MsgChan:
			// 檢查日誌大小
			if isVaild, err := f.checkFileSizeValid(false); err != nil {
				panic(err)
			} else if !isVaild {
				f.divideFile(false)
			}
			// 寫入一般日誌
			fmt.Fprintf(f.FileObj, "[%v] [%v] [%v:%v:%v] %v\n", msg.Time, parseOpenLevelToStr(msg.Level), msg.FileName, msg.FuncName, msg.Line, msg.Msg)

			// error級別以上的要在記錄在errFile
			if msg.Level >= ERROR {
				// 檢查日誌大小
				if isVaild, err := f.checkFileSizeValid(true); err != nil {
					panic(err)
				} else if !isVaild {
					f.divideFile(true)
				}
				fmt.Fprintf(f.ErrorFileObj, "[%v] [%v] [%v:%v:%v] %v\n", msg.Time, parseOpenLevelToStr(msg.Level), msg.FileName, msg.FuncName, msg.Line, msg.Msg)
			}
		default:
			// 沒訊息就先休息一下
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// Debug ...
func (f *FileLogger) Debug(format string, arg ...interface{}) {
	f.log(DEBUG, format, arg...)
}

// Trace ...
func (f *FileLogger) Trace(format string, arg ...interface{}) {
	f.log(TRACE, format, arg...)
}

// Info ...
func (f *FileLogger) Info(format string, arg ...interface{}) {
	f.log(INFO, format, arg...)
}

// Warning ...
func (f *FileLogger) Warning(format string, arg ...interface{}) {
	f.log(WARNING, format, arg...)
}

// Error ...
func (f *FileLogger) Error(format string, arg ...interface{}) {
	f.log(ERROR, format, arg...)
}

// Fatal ...
func (f *FileLogger) Fatal(format string, arg ...interface{}) {
	f.log(FATAL, format, arg...)
}
