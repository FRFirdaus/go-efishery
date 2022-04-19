package logger

import (
	"fmt"

	lr "github.com/sirupsen/logrus"
)

const (
	LevelTraceMsg   string = "[TRACE]"
	LevelDebugMsg   string = "[DEBUG]"
	LevelInfoMsg    string = "[INFO]"
	LevelWarnMsg    string = "[WARNING]"
	LevelErrorMsg   string = "[ERROR]"
	LevelFatalMsg   string = "[FATAL]"
	LevelPanicMsg   string = "[PANIC]"
	LevelUnknownMsg string = "[UNKNOWN LOG LEVEL]"
)

var (
	ErrUnknownLevel error = fmt.Errorf(`[UNKNOWN LOG LEVEL] [FAILED] Logger Error`)
)

func (l *logger) convertAndSetLevel() {
	l.setLevelLogrus()
}

func (l *logger) setLevelLogrus() {
	var lrLevel lr.Level
	switch l.opt.Level {
	case LevelTrace:
		lrLevel = lr.TraceLevel
		l.log.Info(OK, infoLogger, LevelTraceMsg)
	case LevelDebug:
		lrLevel = lr.DebugLevel
		l.log.Info(OK, infoLogger, LevelDebugMsg)
	case LevelInfo:
		lrLevel = lr.InfoLevel
		l.log.Info(OK, infoLogger, LevelInfoMsg)
	case LevelWarn:
		lrLevel = lr.WarnLevel
		l.log.Info(OK, infoLogger, LevelWarnMsg)
	case LevelError:
		lrLevel = lr.ErrorLevel
		l.log.Info(OK, infoLogger, LevelErrorMsg)
	case LevelFatal:
		lrLevel = lr.FatalLevel
		l.log.Info(OK, infoLogger, LevelFatalMsg)
	case LevelPanic:
		lrLevel = lr.PanicLevel
		l.log.Info(OK, infoLogger, LevelPanicMsg)
	default:
		err := ErrUnknownLevel
		l.log.Panic(err)
	}
	//set logrus log level
	l.log.SetLevel(lrLevel)
}
