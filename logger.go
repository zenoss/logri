package logri

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
)

const (
	RootLoggerName              = ""
	NilLevel       logrus.Level = 255
)

var (
	ErrInvalidRootLevel                    = errors.New("The root logger must have a level")
	RootLogger          logrus.FieldLogger = NewRootLogger()
)

type Logger struct {
	mu       *sync.Mutex
	name     string
	parent   *Logger
	absLevel logrus.Level
	children map[string]*Logger
	logger   *logrus.Logger
}

func NewRootLogger() *Logger {
	return &Logger{
		name:     RootLoggerName,
		absLevel: logrus.InfoLevel,
		children: make(map[string]*Logger),
		logger:   logrus.New(),
	}
}

func (l *Logger) GetLogrusLogger() *logrus.Logger {
	return l.logger
}

func (l *Logger) GetChild(name string) *Logger {
	relative := strings.TrimPrefix(name, l.name+".")
	abs := fmt.Sprintf("%s.%s", l.name, relative)
	parent := l
	for _, part := range strings.Split(relative, ".") {
		logger, ok := parent.children[part]
		if !ok {
			logger = &Logger{
				name:     abs,
				parent:   parent,
				absLevel: NilLevel,
				children: make(map[string]*Logger),
				logger: &logrus.Logger{
					Out:       parent.logger.Out,
					Formatter: parent.logger.Formatter,
					Hooks:     parent.logger.Hooks,
					Level:     parent.logger.Level,
				},
			}
			parent.children[part] = logger
		}
		parent = logger
	}
	return parent
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
