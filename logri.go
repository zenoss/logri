package logri

import "github.com/Sirupsen/logrus"

var (
	// RootLogger is the default created logger tree.
	RootLogger logrus.FieldLogger = NewLoggerFromLogrus(logrus.New())
)
