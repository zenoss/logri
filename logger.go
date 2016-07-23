package logri

import (
	"errors"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
)

const (
	RootLoggerName              = ""
	NilLevel       logrus.Level = 255
)

var (
	ErrInvalidRootLevel = errors.New("The root logger must have a level")
)

type Logger struct {
	mu       *sync.Mutex
	name     string
	parent   *Logger
	absLevel logrus.Level
	children []*Logger
	logger   *logrus.Logger
}

func (l *Logger) GetChild(name string) *Logger {
	for _, part := range strings.Split(name, ".") {

	}
	return nil
}

func (l *Logger) SetLevel(level logrus.Level, propagate bool) error {
	if level == l.absLevel {
		return nil
	}
	if level == NilLevel && l.name == RootLoggerName {
		return ErrInvalidRootLevel
	}
	l.absLevel = level
	switch level {
	case NilLevel:
		l.logger.Level = l.parent.GetEffectiveLevel()
	default:
		l.logger.Level = level
	}
	if propagate {
		l.propagate()
	}
	return nil
}

func (l *Logger) GetEffectiveLevel() logrus.Level {
	return l.logger.Level
}

func (l *Logger) GetLevel() logrus.Level {
	return l.absLevel
}

func (l *Logger) propagate() {
	for _, child := range l.children {
		go child.inheritLevel(l.logger.Level)
	}
}

func (l *Logger) inheritLevel(parentLevel logrus.Level) {
	if l.absLevel == NilLevel {
		l.logger.Level = parentLevel
		l.propagate()
	}
}
