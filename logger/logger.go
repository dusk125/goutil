package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	logErr, logWarn, logInfo, logDebug *log.Logger
)

func init() {
	logDebug = log.New(os.Stdout, "[DEBUG] ", 0)
	logInfo = log.New(os.Stdout, "[Info] ", 0)
	logWarn = log.New(os.Stdout, "[Warn] ", 0)
	logErr = log.New(os.Stdout, "[ERROR] ", 0)
}

func SetFlags(flags int) {
	log.SetFlags(flags)
	logDebug.SetFlags(flags)
	logInfo.SetFlags(flags)
	logWarn.SetFlags(flags)
	logErr.SetFlags(flags)
}

func out(logger *log.Logger, v ...interface{}) {
	_ = logger.Output(3, fmt.Sprintln(v...))
}

func outf(logger *log.Logger, f string, v ...interface{}) {
	if !strings.HasSuffix(f, "\n") {
		f += "\n"
	}
	_ = logger.Output(3, fmt.Sprintf(f, v...))
}

func Debug(v ...interface{}) {
	out(logDebug, v...)
}

func Debugf(f string, v ...interface{}) {
	outf(logDebug, f, v...)
}

func Info(v ...interface{}) {
	out(logInfo, v...)
}

func Infof(f string, v ...interface{}) {
	outf(logInfo, f, v...)
}

func Warn(v ...interface{}) {
	out(logWarn, v...)
}

func Warnf(f string, v ...interface{}) {
	outf(logWarn, f, v...)
}

func Error(v ...interface{}) {
	out(logErr, v...)
}

func Errorf(f string, v ...interface{}) {
	outf(logErr, f, v...)
}
