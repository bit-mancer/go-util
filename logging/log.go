/*
Package logging contains a simple wrapper around a concrete structured logger.
*/
package logging

import (
	"os"

	"github.com/Sirupsen/logrus"
)

// SourceKey is the map key used for logging sources (see NewLog).
const SourceKey = "_source"

// DomainKey is the map key used for logging domains (see NewDomainLogger).
const DomainKey = "_domain"

// HostKey is the map key used for logging the machine hostname (via os.Hostname).
const HostKey = "_host"

// AuditKey is the map key used for audit-related entries (see Audit).
const AuditKey = "_audit"

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
// 	* 'source' is the optional source of the events (the root that all logs under this logger and its domains will
//    be tagged with).
func NewLog(source string) Logger {
	baseFormatter := &logrus.TextFormatter{}
	baseFormatter.FullTimestamp = true

	utcFormatter := &utcFormatter{baseFormatter}

	l := logrus.New()
	l.Formatter = utcFormatter
	l.Out = os.Stdout
	l.Level = logrus.DebugLevel

	fields := logrus.Fields{}

	if source != "" {
		fields[SourceKey] = source
	}

	if hostname, err := os.Hostname(); err == nil {
		fields[HostKey] = hostname
	}

	return (*logger)(l.WithFields(fields))
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

// Audit logs critical auditing-related information -- this level is suitable for production; the log should be
// *actionable information* that will be read by a human, or by a machine.
// Log entries made with Audit will be tagged with the audit field (AuditKey: true).
func (l *logger) Audit(args ...interface{}) {
	(*logrus.Entry)(l).WithField(AuditKey, true).Warn(args)
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
	e := (*logrus.Entry)(l).WithField(DomainKey, domain)
	return (*logger)(e)
}
