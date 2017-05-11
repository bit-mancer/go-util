package logging

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SourceKey is the map key used for logging sources (see NewLog).
const SourceKey = "_source"

// DomainKey is the map key used for logging domains (see NewDomainLogger).
const DomainKey = "_domain"

// HostKey is the map key used for logging the machine hostname (via os.Hostname).
const HostKey = "_host"

// FieldSet is an alias for a slice of zapcore Fields. Functions that aggregate fields can return a FieldSet to allow
// further aggregation.
type FieldSet []zapcore.Field

// AppendFieldSet appends a FieldSet to the current FieldSet.
func (f FieldSet) AppendFieldSet(fields FieldSet) FieldSet {
	return append([]zapcore.Field(f), fields...)
}

// Append appends fields to the current FieldSet.
func (f FieldSet) Append(fields ...zapcore.Field) FieldSet {
	return append([]zapcore.Field(f), fields...)
}

// Logger is implemented by... loggers.
type Logger interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	ErrorWithTrace(err error, msg string, fields ...zapcore.Field)

	NewDomainLogger(domain string) Logger
}

type logger zap.Logger

// NewLog returns a new Logger instance.
// If 'productionLogging' is true, JSON logger is constructed which logs at the Info level (by default); if
// 'productionLogging' is false, a developent-oriented console logger is constructed which logs at the Debug level.
// 'source' is the optional source (e.g. application or service name) of the events; all logs under this logger and
// its domains will be tagged with the source, if one is provided.
func NewLog(productionLogging bool, source string) (Logger, error) {

	var log *zap.Logger
	var err error

	if productionLogging {
		log, err = zap.NewProduction()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create production logger: %v\n", err)
			os.Exit(1)
		}
	} else {
		log, err = zap.NewDevelopment()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create development logger: %v\n", err)
			os.Exit(1)
		}
	}

	if source != "" {
		log = log.With(zap.String(SourceKey, source))
	}

	if hostname, err := os.Hostname(); err == nil {
		log = log.With(zap.String(HostKey, hostname))
	}

	return (*logger)(log), nil
}

func (l *logger) Debug(msg string, fields ...zapcore.Field) {
	(*zap.Logger)(l).Debug(msg, fields...)
}

func (l *logger) Info(msg string, fields ...zapcore.Field) {
	(*zap.Logger)(l).Info(msg, fields...)
}

func (l *logger) ErrorWithTrace(err error, msg string, fields ...zapcore.Field) {
	fields = append(fields, zap.Error(err))
	(*zap.Logger)(l).Error(msg, fields...)
}

func (l *logger) NewDomainLogger(domain string) Logger {
	newLogger := (*zap.Logger)(l).With(zap.String(DomainKey, domain))
	return (*logger)(newLogger)
}
