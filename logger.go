package logri

import "github.com/Sirupsen/logrus"

type Logger struct {
	name     string
	parent   *Logger
	level    logrus.Level
	children []*Logger
	logrus.Logger
}

func (l *Logger) SetLevel(level logrus.Level) {
	l.Level = level
}

func (l *Logger) GetEffectiveLevel() logrus.Level {
	return l.Level
}

func (l *Logger) GetLevel() logrus.Level {
	return l.Level
}
