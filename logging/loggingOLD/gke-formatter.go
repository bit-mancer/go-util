package loggingOLD

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pquerna/ffjson/ffjson"
)

// gkeFormatter is a log formatter that produces structured logs compatible with Google Container Engine (Stackdriver
// via fluentd).
type gkeFormatter struct {
}

func (f *gkeFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	// Golang reference layout:
	// Mon Jan 2 15:04:05 MST 2006
	// 01/02 03:04:05PM '06 -0700
	const timeLayout = "2006-01-02T15:04:05.000Z"

	msg := gkeMessage{
		Time:     entry.Time.UTC().Format(timeLayout),
		Severity: strings.ToUpper(entry.Level.String()),
		Message:  entry.Message,
	}

	if len(entry.Data) > 0 {
		msg.Meta = entry.Data
	}

	data, err := ffjson.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal GKE log message: %v", err)
	}
	return data, nil
}
