package logger

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"sort"
	"strings"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/sirupsen/logrus"
)

// LogLevel represents log level type
type LogLevel string

// Log level constants
const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
	LevelFatal LogLevel = "fatal"
)

// ANSI color codes
const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
	colorReset  = "\033[0m"
)

type CustomFormatter struct {
	ForceColor bool // Whether to force color, even in non-terminal environments
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")
	level := strings.ToUpper(entry.Level.String())

	// Set color based on log level
	var levelColor, resetColor string
	if f.ForceColor {
		switch entry.Level {
		case logrus.DebugLevel:
			levelColor = colorCyan
		case logrus.InfoLevel:
			levelColor = colorGreen
		case logrus.WarnLevel:
			levelColor = colorYellow
		case logrus.ErrorLevel:
			levelColor = colorRed
		case logrus.FatalLevel:
			levelColor = colorPurple
		default:
			levelColor = colorReset
		}
		resetColor = colorReset
	}

	// Extract caller field
	caller := ""
	if val, ok := entry.Data["caller"]; ok {
		caller = fmt.Sprintf("%v", val)
	}

	// Concatenate field part: other fields sorted and output (request_id excluded)
	fields := ""

	// Fields sorted and output (request_id and caller excluded)
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		if k != "caller" && k != "request_id" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		if f.ForceColor {
			val := fmt.Sprintf("%v", entry.Data[k])
			coloredVal := fmt.Sprintf("%s%s%s", colorWhite, val, colorReset)
			if k == "error" {
				coloredVal = fmt.Sprintf("%s%s%s", colorRed, val, colorReset)
			}
			fields += fmt.Sprintf("%s%s%s=%s ",
				colorCyan, k, colorReset, coloredVal)
		} else {
			fields += fmt.Sprintf("%s=%v ", k, entry.Data[k])
		}
	}

	fields = strings.TrimSpace(fields)

	// Concatenate final output content, add color
	if f.ForceColor {
		coloredTimestamp := fmt.Sprintf("%s%s%s", colorGray, timestamp, resetColor)
		coloredCaller := caller
		if caller != "" {
			coloredCaller = fmt.Sprintf("%s%s%s", colorPurple, caller, resetColor)
		}
		return []byte(fmt.Sprintf("%s%-5s%s[%s] [%s] %-20s | %s\n",
			levelColor, level, resetColor, coloredTimestamp, fields, coloredCaller, entry.Message)), nil
	}

	return []byte(fmt.Sprintf("%-5s[%s] [%s] %-20s | %s\n",
		level, timestamp, fields, caller, entry.Message)), nil
}

// Initialize global log settings
func init() {
	// Set log format without modifying global timezone
	logrus.SetFormatter(&CustomFormatter{ForceColor: true})
	logrus.SetReportCaller(false)
}

// GetLogger gets logger instance
func GetLogger(c context.Context) *logrus.Entry {
	if logger := c.Value(types.LoggerContextKey); logger != nil {
		return logger.(*logrus.Entry)
	}
	newLogger := logrus.New()
	newLogger.SetFormatter(&CustomFormatter{ForceColor: true})
	// Set default log level
	newLogger.SetLevel(logrus.DebugLevel)
	// Enable caller information
	return logrus.NewEntry(newLogger)
}

// SetLogLevel sets log level
func SetLogLevel(level LogLevel) {
	var logLevel logrus.Level

	switch level {
	case LevelDebug:
		logLevel = logrus.DebugLevel
	case LevelInfo:
		logLevel = logrus.InfoLevel
	case LevelWarn:
		logLevel = logrus.WarnLevel
	case LevelError:
		logLevel = logrus.ErrorLevel
	case LevelFatal:
		logLevel = logrus.FatalLevel
	default:
		logLevel = logrus.InfoLevel
	}

	logrus.SetLevel(logLevel)
}

// Add caller field
func addCaller(entry *logrus.Entry, skip int) *logrus.Entry {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return entry
	}
	shortFile := path.Base(file)
	funcName := "unknown"
	if fn := runtime.FuncForPC(pc); fn != nil {
		// Keep only function name, without package path (e.g., doSomething)
		fullName := path.Base(fn.Name())
		parts := strings.Split(fullName, ".")
		funcName = parts[len(parts)-1]
	}
	return entry.WithField("caller", fmt.Sprintf("%s:%d[%s]", shortFile, line, funcName))
}

// WithRequestID adds request ID to log
func WithRequestID(c context.Context, requestID string) context.Context {
	return WithField(c, "request_id", requestID)
}

// WithField adds a field to log
func WithField(c context.Context, key string, value interface{}) context.Context {
	logger := GetLogger(c).WithField(key, value)
	return context.WithValue(c, types.LoggerContextKey, logger)
}

// WithFields adds multiple fields to log
func WithFields(c context.Context, fields logrus.Fields) context.Context {
	logger := GetLogger(c).WithFields(fields)
	return context.WithValue(c, types.LoggerContextKey, logger)
}

// Debug outputs debug level log
func Debug(c context.Context, args ...interface{}) {
	addCaller(GetLogger(c), 2).Debug(args...)
}

// Debugf outputs debug level log using formatted string
func Debugf(c context.Context, format string, args ...interface{}) {
	addCaller(GetLogger(c), 2).Debugf(format, args...)
}

// Info outputs info level log
func Info(c context.Context, args ...interface{}) {
	addCaller(GetLogger(c), 2).Info(args...)
}

// Infof outputs info level log using formatted string
func Infof(c context.Context, format string, args ...interface{}) {
	addCaller(GetLogger(c), 2).Infof(format, args...)
}

// Warn outputs warning level log
func Warn(c context.Context, args ...interface{}) {
	addCaller(GetLogger(c), 2).Warn(args...)
}

// Warnf outputs warning level log using formatted string
func Warnf(c context.Context, format string, args ...interface{}) {
	addCaller(GetLogger(c), 2).Warnf(format, args...)
}

// Error outputs error level log
func Error(c context.Context, args ...interface{}) {
	addCaller(GetLogger(c), 2).Error(args...)
}

// Errorf outputs error level log using formatted string
func Errorf(c context.Context, format string, args ...interface{}) {
	addCaller(GetLogger(c), 2).Errorf(format, args...)
}

// ErrorWithFields outputs error level log with additional fields
func ErrorWithFields(c context.Context, err error, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	addCaller(GetLogger(c), 2).WithFields(fields).Error("发生错误")
}

// Fatal 输出致命级别的日志并退出程序
func Fatal(c context.Context, args ...interface{}) {
	addCaller(GetLogger(c), 2).Fatal(args...)
}

// Fatalf 使用格式化字符串输出致命级别的日志并退出程序
func Fatalf(c context.Context, format string, args ...interface{}) {
	addCaller(GetLogger(c), 2).Fatalf(format, args...)
}

// CloneContext 复制上下文中的关键信息到新上下文
func CloneContext(ctx context.Context) context.Context {
	newCtx := context.Background()

	for _, k := range []types.ContextKey{
		types.LoggerContextKey,
		types.TenantIDContextKey,
		types.RequestIDContextKey,
		types.TenantInfoContextKey,
	} {
		if v := ctx.Value(k); v != nil {
			newCtx = context.WithValue(newCtx, k, v)
		}
	}

	return newCtx
}
