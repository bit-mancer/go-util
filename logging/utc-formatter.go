package logging

import (
	"github.com/Sirupsen/logrus"
)

// utcFormatter is a logrus formatter that logs timestamps in UTC.
type utcFormatter struct {
	baseFormatter logrus.Formatter
}

func (f *utcFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return f.baseFormatter.Format(e)
}
