package logger

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	OutputStdout  string = `stdout`
	OutputDiscard string = `discard`
	OutputCustom  string = `custom`

	outputStdout  string = `[STDOUT]`
	outputDiscard string = `[DISCARD]`
	outputCustom  string = `[CUSTOM]`
	outputUnknown string = `[UNKNOWN LOG OUTPUT]`
)

var (
	ErrUnknownOutput = fmt.Errorf(`[UNKNOWN LOG OUTPUT] [FAILED] Logger Error`)
)

func (l *logger) convertAndSetOutput() {
	switch l.opt.Output {
	case OutputDiscard:
		l.log.Info(OK, infoLogger, outputDiscard)
		l.log.SetOutput(ioutil.Discard)
	case OutputStdout:
		l.log.SetOutput(os.Stdout)
		l.log.Info(OK, infoLogger, outputStdout)
	case OutputCustom:
		l.log.SetOutput(l.opt.CustomWriter)
		l.log.Info(OK, infoLogger, outputCustom)
	default:
		l.log.Panic(ErrUnknownOutput)
	}
}
