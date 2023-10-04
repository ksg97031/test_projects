package logging

import (
	"fmt"
	"os"

	"github.com/rifflock/lfshook"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var stdFormatter *TextFormatter  // Command Line Output Format
var fileFormatter *TextFormatter // File output format

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.SetReportCaller(true)
	stdFormatter = &TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02.15:04:05",
		ForceFormatting: true,
		ForceColors:     true,
		DisableColors:   false,
		ReportCaller:    true,
	}
	fileFormatter = &TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02.15:04:05",
		ForceFormatting: true,
		ForceColors:     false,
		DisableColors:   true,
		ReportCaller:    true,
	}

	Logger.SetFormatter(stdFormatter)
	Logger.SetLevel(logrus.DebugLevel)

	logPath, _ := os.Getwd()
	logName := fmt.Sprintf("%s/logs/yi_log_", logPath)
	writer, _ := rotatelogs.New(logName + "%Y_%m_%d" + ".log")
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.DebugLevel: writer,
		logrus.ErrorLevel: writer,
	}, fileFormatter)
	Logger.SetOutput(os.Stdout)
	Logger.AddHook(lfHook)
}
