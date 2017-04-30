/*
Package logging contains a simple wrapper around a concrete structured logger.
*/
package logging

import (
	"os"

	"github.com/Sirupsen/logrus"
)

const sourceKey = "_source"
const domainKey = "_domain"
const auditKey = "_audit"

// Fields is a set of strutured log information (i.e. a map-compatible type).
// Example: log.WithFields(logging.Fields{"myKey": myVal, "myOtherKey": myOtherVal})
type Fields logrus.Fields

// Logger is implemented by structured loggers.
type Logger interface {
	WithField(key string, value interface{}) *logger
	WithFields(fields Fields) *logger
	WithError(err error) *logger

	Info(args ...interface{})
	Debug(args ...interface{})
	Audit(args ...interface{})

	SetEnableDebug(bool)

	NewDomainLogger(domain string) *logger
}

type logger logrus.Entry

// NewLog returns a new Logger instance.
// 	* 'source' is the optional source of the events (the root process that all logs under this logger and its domains
//    will be tagged with).
func NewLog(source string) Logger {
	baseFormatter := &logrus.TextFormatter{}
	baseFormatter.FullTimestamp = true

	utcFormatter := &utcFormatter{baseFormatter}

	l := logrus.New()
	l.Formatter = utcFormatter
	l.Out = os.Stdout
	l.Level = logrus.DebugLevel

	if source != "" {
		return (*logger)(l.WithField(sourceKey, source))
	}
	return (*logger)(l.WithFields(logrus.Fields{}))
}

// WithField logs the provided key and value as a structured log entry, and returns the entry for further use.
func (l *logger) WithField(key string, value interface{}) *logger {
	e := (*logrus.Entry)(l).WithField(key, value)
	return (*logger)(e)
}

// WithFields logs the provided fields as a structured log entry, and returns the entry for further use.
func (l *logger) WithFields(fields Fields) *logger {
	e := (*logrus.Entry)(l).WithFields(logrus.Fields(fields))
	return (*logger)(e)
}

// WithError logs the provided error as a structured log entry, and returns the entry for further use.
func (l *logger) WithError(err error) *logger {
	e := (*logrus.Entry)(l).WithError(err)
	return (*logger)(e)
}

// Debug logs development-related information that is not typically necessary for production. When choosing a
// logging level, you should assume that debug logging will be disabled in production environments.
func (l *logger) Debug(args ...interface{}) {
	(*logrus.Entry)(l).Debug(args)
}

// Info logs at a level suitable for production; the log should be *actionable information* that will be read by a
// human, or by a machine.
func (l *logger) Info(args ...interface{}) {
	(*logrus.Entry)(l).Info(args)
}

// Info logs at the "warn" level -- this level is suitable for production; the log should be *actionable information*
// that will be read by a human, or by a machine.
func (l *logger) Audit(args ...interface{}) {
	(*logrus.Entry)(l).WithField(auditKey, true).Warn(args)
}

// SetEnableDebug configures whether the "debug" log level is enabled (the default is true).
func (l *logger) SetEnableDebug(enable bool) {
	if enable {
		l.Logger.Level = logrus.DebugLevel
		l.Logger.Warn("Log level set to DEBUG")
	} else {
		l.Logger.Level = logrus.InfoLevel
		l.Logger.Warn("Log level set to INFO")
	}
}

// NewDomainLogger returns a new Logger instance with the provided domain.
// The domain is set via WithField() with a well-defined key.
// The new Logger will take on all attributes of the current Logger.
func (l *logger) NewDomainLogger(domain string) *logger {
	e := (*logrus.Entry)(l).WithField(domainKey, domain)
	return (*logger)(e)
}
