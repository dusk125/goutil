package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

type Level uint32

const (
	LevelFatal Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

func (l Level) String() string {
	switch l {
	case LevelFatal:
		return "Fatal"
	case LevelError:
		return "Error"
	case LevelWarn:
		return "Warn"
	case LevelInfo:
		return "Info"
	case LevelDebug:
		return "Debug"
	case LevelTrace:
		return "Trace"
	}
	return ""
}

var (
	loggers = [LevelTrace + 1]*log.Logger{}
	level   uint32
)

func init() {
	l := LevelInfo
	if level, has := os.LookupEnv("LOG_LEVEL"); has {
		switch strings.ToLower(level) {
		case "fatal":
			l = LevelFatal
		case "error":
			l = LevelError
		case "warn":
			l = LevelWarn
		case "info":
			l = LevelInfo
		case "debug":
			l = LevelDebug
		case "trace":
			l = LevelTrace
		}
	}
	SetLevel(l)
	for i := range loggers {
		loggers[i] = log.New(os.Stdout, fmt.Sprintf("[%v] ", Level(i)), 0)
	}
}

// Sets the given flags to all of the loggers, including the default package 'log' logger.
func SetFlags(flags int) {
	log.SetFlags(flags)
	for _, logger := range loggers {
		logger.SetFlags(flags)
	}
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
	for _, logger := range loggers {
		logger.SetOutput(w)
	}
}

// Sets the filtering level for the logging system.
//
//	A value will allow all message at, and below, it's logging level.
//	For example, setting the logging level to LevelDebug will allow all Debug, Info, Warn, Error, and Fatal messages through.
//	Setting the level to LevelError will allow only Error and Fatal through.
func SetLevel(l Level) {
	atomic.StoreUint32(&level, uint32(l))
}

func out(logger Level, v ...any) {
	if Level(atomic.LoadUint32(&level)) >= logger {
		_ = loggers[logger].Output(3, fmt.Sprintln(v...))
	}
}

func outf(logger Level, f string, v ...any) {
	if Level(atomic.LoadUint32(&level)) >= logger {
		if !strings.HasSuffix(f, "\n") {
			f += "\n"
		}
		_ = loggers[logger].Output(3, fmt.Sprintf(f, v...))
	}
}

func Trace(v ...any) {
	out(LevelTrace, v...)
}

func Tracef(f string, v ...any) {
	outf(LevelTrace, f, v...)
}

func Debug(v ...any) {
	out(LevelDebug, v...)
}

func Debugf(f string, v ...any) {
	outf(LevelDebug, f, v...)
}

func Info(v ...any) {
	out(LevelInfo, v...)
}

func Infof(f string, v ...any) {
	outf(LevelInfo, f, v...)
}

func Warn(v ...any) {
	out(LevelWarn, v...)
}

func Warnf(f string, v ...any) {
	outf(LevelWarn, f, v...)
}

func Error(v ...any) {
	out(LevelError, v...)
}

func Errorf(f string, v ...any) {
	outf(LevelError, f, v...)
}

func Fatal(v ...any) {
	out(LevelFatal, v...)
	os.Exit(1)
}

func Fatalf(f string, v ...any) {
	outf(LevelFatal, f, v...)
	os.Exit(1)
}

func ErrorWriter() io.Writer {
	return loggers[LevelError].Writer()
}
