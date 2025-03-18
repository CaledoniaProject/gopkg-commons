package commons

import (
	"fmt"
	"os"
	"path"
	"runtime"

	lumberjack "github.com/CaledoniaProject/gopkg-lumberjack"
	"github.com/sirupsen/logrus"
)

func SetupConsoleLogger() {
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(GetTTYFormatter())
}

func SetupJsonFileLogger(filename string) {
	_, l_logger := GetRotatingJSONLogger(filename)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(l_logger)
}

func GetCleanJSONFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", ""
		},
	}
}

func GetTextFormatter() *logrus.TextFormatter {
	return &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   true,
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("[%s:%d]", path.Base(f.File), f.Line)
		},
	}
}

func GetTTYFormatter() *logrus.TextFormatter {
	return &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
		ForceColors:     true,
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("[%s:%d]", path.Base(f.File), f.Line)
		},
	}
}

func GetRotatingFileLogger(filename string) (*logrus.Logger, *lumberjack.Logger) {
	lumberjack_logger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    200,
		MaxBackups: 20,
		MaxAge:     180,
		Compress:   true,
	}
	logger := &logrus.Logger{
		Out:          lumberjack_logger,
		Formatter:    &logrus.TextFormatter{},
		ReportCaller: true,
		Level:        logrus.InfoLevel,
	}

	return logger, lumberjack_logger
}

func GetRotatingJSONLogger(filename string) (*logrus.Logger, *lumberjack.Logger) {
	lumberjack_logger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    200,
		MaxBackups: 20,
		MaxAge:     180,
		Compress:   true,
	}
	logger := &logrus.Logger{
		Out:          lumberjack_logger,
		Formatter:    &logrus.JSONFormatter{},
		ReportCaller: true,
		Level:        logrus.InfoLevel,
	}

	return logger, lumberjack_logger
}
