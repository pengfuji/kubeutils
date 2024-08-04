/*
 * @Time : 2024/7/31 18:12
 * @Author : diehao.yuan
 * @Email : diehao.yuan@outlook.com
 * @File : log.go
 */
package log

import (
	"fmt"
	"path"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// CustomJSONFormatter 是一个自定义的JSON日志格式化器，用于指定日志的输出格式。
type CustomJSONFormatter struct {
	logrus.JSONFormatter
}

// Format 方法实现了logrus.Formatter接口，用于自定义日志的JSON格式。
func (f *CustomJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 自定义日志格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := entry.Level.String()
	msg := entry.Message
	funcVal := fmt.Sprintf("%s:%d %s", entry.Caller.File, entry.Caller.Line, entry.Caller.Function)
	jsonStr := fmt.Sprintf(`{"timestamp": "%s", "level": "%s", "func": "%s", "file": "%s", "msg": "%s"}`, timestamp, level, funcVal, entry.Caller.File, msg)
	// 添加换行符以便于阅读日志文件
	jsonStr += "\n"
	return []byte(jsonStr), nil
}

// InitLogger 初始化日志系统，配置日志级别、滚动策略等
func InitLogger(logLevel string, logPath string, maxBackups, maxAge int, compress bool) {
	logFileNamePattern := "app-%s.log"

	// 构造日志文件的完整路径与名称，使用当前日期格式化文件名
	currentDate := time.Now().Format("20060102")
	logFilePath := path.Join(logPath, fmt.Sprintf(logFileNamePattern, currentDate))

	// 配置日志滚动策略,设置日志输出到lumberjack，实现日志滚动和压缩
	logWriter := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxBackups: maxBackups, // 设置最大备份文件数
		MaxAge:     maxAge,     // 设置文件保留的最大天数
		Compress:   compress,   // 设置是否启用压缩禁用的日志文件
		LocalTime:  true,       // 设置是否使用本地时间进行文件滚动
	}

	// 设置logrus的日志输出目标为上面配置的滚动日志写入器
	logrus.SetOutput(logWriter)

	// 设置自定义的JSON格式
	logrus.SetFormatter(&CustomJSONFormatter{})

	// 解析并设置日志级别，如果配置文件中的级别不合法，则使用默认的Info界别
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel // 设置默认级别为info
		fmt.Printf("警告：配置文件中的日志级别不合法，已设置为默认级别：Info\n")
	}
	logrus.SetLevel(level)

	// 启用在日志中添加调用者信息（文件名、函数名和行号）
	logrus.SetReportCaller(true)
}

// Debug 输出Debug级别的日志信息
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Debugf 格式化输出Debug级别的日志信息
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Info 输出Info级别的日志信息
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Infof 格式化输出Info级别的日志信息
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Warn 输出Warning级别的日志信息
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

// Warnf 格式化输出Warning级别的日志信息
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

// Error 输出Error级别的日志信息
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Errorf 格式化输出Error级别的日志信息
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}
