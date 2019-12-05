package common

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// log
var Logger *logrus.Logger

func loggerInit() {
	Logger = logrus.New()
	Logger.AddHook(NewFileHook())
}

type FileHook struct {
	accessPath string
	errPath    string
	formatter  logrus.Formatter
	lock       *sync.Mutex // 日志锁 防止 并发写
}

// 默认formatter
var defaultFormatter = &logrus.JSONFormatter{}

// 默认文件名
var defaultAccessName string = "log_access_" + time.Now().Format("2006-01-02") + ".log"

// 默认错误文件名
var defaultErrName string = "log_error_" + time.Now().Format("2006-01-02") + ".log"

// 默认路径
var defaultPath = "/log/"

// new hook
func NewFileHook() *FileHook {

	defaultBasePath, err := os.Getwd()
	if err != nil {
		log.Fatal("common/logger # NewFileHook os.GetWd error :", err)
	}
	var defaultAccessPath string = defaultBasePath + defaultPath + defaultAccessName
	var defaultErrPath string = defaultBasePath + defaultPath + defaultErrName

	return New(defaultAccessPath, defaultErrPath, defaultFormatter)
}

// New
func New(accessPath string, errPath string, formatter logrus.Formatter) *FileHook {

	fileHook := &FileHook{
		lock: new(sync.Mutex),
	}

	fileHook.SetFormatter(formatter)
	fileHook.SetAccessPath(accessPath)
	fileHook.SetErrPath(errPath)
	return fileHook
}

// 适用的日志级别
func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// hook接口的调用逻辑
func (hook *FileHook) Fire(entry *logrus.Entry) error {

	// error 写到error文件
	// 正常 写到普通文件
	switch entry.Level {
	case logrus.PanicLevel:
		fallthrough
	case logrus.FatalLevel:
		fallthrough
	case logrus.ErrorLevel:
		return hook.fileWrite(entry, hook.errPath)
	case logrus.WarnLevel:
		fallthrough
	case logrus.InfoLevel:
		fallthrough
	case logrus.DebugLevel, logrus.TraceLevel:
		fallthrough
	default:
		return hook.fileWrite(entry, hook.accessPath)
	}
}

// 设置formatter
func (hook *FileHook) SetFormatter(formatter logrus.Formatter) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if formatter == nil {
		formatter = defaultFormatter
	} else {
		switch formatter.(type) {
		case *logrus.TextFormatter:
			textFormatter := formatter.(*logrus.TextFormatter)
			textFormatter.DisableColors = true
		case *logrus.JSONFormatter:
			// do nothing
		}

	}

	hook.formatter = formatter
}

/**
设置 日志路径
*/
func (hook *FileHook) SetAccessPath(path string) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.accessPath = path
}

/**
设置 错误日志路径
*/
func (hook *FileHook) SetErrPath(path string) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.errPath = path
}

/**
写入文件
*/
func (hook *FileHook) fileWrite(entry *logrus.Entry, path string) error {
	var (
		fd  *os.File
		err error
	)

	// 加锁 防止并发写入
	hook.lock.Lock()
	defer hook.lock.Unlock()

	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("path :", path)
	fd, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("failed to open logfile:", path, err)
		return err
	}
	defer fd.Close()

	content, err := hook.formatter.Format(entry)

	if err != nil {
		log.Fatal("failed to generate string for entry:", err)
		return err
	}
	_, err = fd.Write(content)
	if err != nil {
		log.Fatal("common/logger # fileWrite error :", err)
	}

	return nil
}
