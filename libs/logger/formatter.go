package logger

import (
	"fmt"

	lr "github.com/sirupsen/logrus"
)

const (
	FormatJSON string = "json"
	FormatText string = "text"

	formatJSON string = "[JSON]"
	formatText string = "[TEXT]"
)

var (
	ErrUnknownFormat = fmt.Errorf(`[UNKNOWN LOG FORMAT] [FAILED] Logger Error`)
)

func (l *logger) convertAndSetFormatter() {
	switch l.opt.Formatter {
	case FormatText:
		l.log.SetFormatter(&lr.TextFormatter{})
		l.log.Info(OK, infoLogger, formatText)
	case FormatJSON:
		l.log.SetFormatter(&lr.JSONFormatter{})
		l.log.Info(OK, infoLogger, formatJSON)
	default:
		l.log.Panic(ErrUnknownFormat)
	}
}
