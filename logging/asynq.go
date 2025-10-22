package logging

import "go.uber.org/zap"

type AsynqLogger struct {
	log *zap.SugaredLogger
}

func NewAsynqLogger(base *zap.Logger) *AsynqLogger {
	return &AsynqLogger{log: base.Sugar()}
}

func (l *AsynqLogger) Debug(args ...any) {
	l.log.Debug(args...)
}

func (l *AsynqLogger) Info(args ...any) {
	l.log.Info(args...)
}

func (l *AsynqLogger) Warn(args ...any) {
	l.log.Warn(args...)
}

func (l *AsynqLogger) Error(args ...any) {
	l.log.Error(args...)
}

func (l *AsynqLogger) Fatal(args ...any) {
	l.log.Fatal(args...)
}
