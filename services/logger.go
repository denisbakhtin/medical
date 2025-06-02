package services

import "log"

// Logger is an interface for web-site logging service
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

// StdLogger implements Logger interface, uses standard "log" package
type StdLogger struct {
	logger *log.Logger
}

// NewStdLogger returns a new instance of stdLogger
func NewStdLogger() *StdLogger {
	return &StdLogger{logger: log.Default()}
}

// Error logs an error, arguments are similar to log.Print
func (l *StdLogger) Error(v ...any) {
	args := append([]any{errorPrefix}, v...)
	l.logger.Print(args...)
}

// Errorf logs a formatted error, arguments are similar to log.Printf
func (l *StdLogger) Errorf(format string, v ...any) {
	l.logger.Printf(errorPrefix+format, v...)
}

// Warning logs a warning, arguments are similar to log.Print
func (l *StdLogger) Warning(v ...any) {
	args := append([]any{warnPrefix}, v...)
	l.logger.Print(args...)
}

// Warningf logs a formatted warning, arguments are similar to log.Printf
func (l *StdLogger) Warningf(format string, v ...any) {
	l.logger.Printf(warnPrefix+format, v...)
}

// Info logs an info message, arguments are similar to log.Print
func (l *StdLogger) Info(v ...any) {
	args := append([]any{infoPrefix}, v...)
	l.logger.Print(args...)
}

// Infof logs a formatted info message, arguments are similar to log.Printf
func (l *StdLogger) Infof(format string, v ...any) {
	l.logger.Printf(infoPrefix+format, v...)
}

// Fatal logs a fatal message with subsequent program termination, arguments are similar to log.Fatal
func (l *StdLogger) Fatal(v ...any) {
	args := append([]any{fatalPrefix}, v)
	l.logger.Fatal(args...)
}
