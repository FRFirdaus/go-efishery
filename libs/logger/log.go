package logger

import (
	"context"

	lr "github.com/sirupsen/logrus"
)

const (
	LevelTrace string = "trace"
	LevelDebug string = "debug"
	LevelInfo  string = "info"
	LevelWarn  string = "warn"
	LevelError string = "error"
	LevelFatal string = "fatal"
	LevelPanic string = "panic"
)

func (l *logger) parseContextFields(ctx context.Context) *lr.Entry {
	doLog := l.logEntry
	if ctx != nil {
		for k, v := range l.opt.ContextFields {
			if val := ctx.Value(v); val != nil {
				doLog = doLog.WithField(k, val)
			}
		}
	}
	return doLog
}

func (l *logger) TraceWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Trace(v...)
}

func (l *logger) Trace(v ...interface{}) {
	l.TraceWithContext(context.TODO(), v...)
}

func (l *logger) DebugWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Debug(v...)
}

func (l *logger) Debug(v ...interface{}) {
	l.DebugWithContext(context.TODO(), v...)
}

func (l *logger) InfoWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Info(v...)
}

func (l *logger) Info(v ...interface{}) {
	l.InfoWithContext(context.TODO(), v...)
}

func (l *logger) WarnWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Warn(v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.WarnWithContext(context.TODO(), v...)
}

func (l *logger) ErrorWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Error(v...)
}

func (l *logger) Error(v ...interface{}) {
	l.ErrorWithContext(context.TODO(), v...)
}

func (l *logger) FatalWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Fatal(v...)
}

func (l *logger) Fatal(v ...interface{}) {
	l.FatalWithContext(context.TODO(), v...)
}

func (l *logger) PanicWithContext(ctx context.Context, v ...interface{}) {
	l.parseContextFields(ctx).Panic(v...)
}

func (l *logger) Panic(v ...interface{}) {
	l.PanicWithContext(context.TODO(), v...)
}
