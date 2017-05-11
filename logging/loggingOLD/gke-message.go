package loggingOLD

import "github.com/Sirupsen/logrus"

//go:generate ffjson $GOFILE

// gkeMessage is a log format that Google Container Engine (Stackdriver via fluentd) understands.
type gkeMessage struct {
	Time     string        `json:"time,omitempty"`
	Severity string        `json:"severity,omitempty"`
	Message  string        `json:"message,omitempty"`
	Meta     logrus.Fields `json:"meta,omitempty"`
}
