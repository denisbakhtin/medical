package services

import "log"

type Logger interface {
	Error(v ...any)
	Errorf(format string, v ...any)
	Warning(v ...any)
	Warningf(format string, v ...any)
	Info(v ...any)
	Infof(format string, v ...any)
	Fatal(v ...any)
}

const errorPrefix = "ERROR: "
const warnPrefix = "WARNING: "
const infoPrefix = "INFO: "
const fatalPrefix = "FATAL: "

type stdLogger struct {
	logger *log.Logger
}

func NewStdLogger() *stdLogger {
	return &stdLogger{logger: log.Default()}
}

func (l *stdLogger) Error(v ...any) {
	args := append([]any{errorPrefix}, v...)
	l.logger.Print(args...)
}

func (l *stdLogger) Errorf(format string, v ...any) {
	l.logger.Printf(errorPrefix+format, v...)
}

func (l *stdLogger) Warning(v ...any) {
	args := append([]any{warnPrefix}, v...)
	l.logger.Print(args...)
}

func (l *stdLogger) Warningf(format string, v ...any) {
	l.logger.Printf(warnPrefix+format, v...)
}

func (l *stdLogger) Info(v ...any) {
	args := append([]any{infoPrefix}, v...)
	l.logger.Print(args...)
}

func (l *stdLogger) Infof(format string, v ...any) {
	l.logger.Printf(infoPrefix+format, v...)
}

func (l *stdLogger) Fatal(v ...any) {
	args := append([]any{fatalPrefix}, v)
	l.logger.Fatal(args...)
}
