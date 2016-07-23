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
	Name     string
	parent   *Logger
	absLevel logrus.Level
	inherit  bool
	children map[string]*Logger
	logger   *logrus.Logger
}

func NewRootLogger() *Logger {
	return NewRootLoggerFromLogrus(logrus.New())
}

func NewRootLoggerFromLogrus(base *logrus.Logger) *Logger {
	return &Logger{
		Name:     RootLoggerName,
		absLevel: logrus.InfoLevel,
		inherit:  true,
		children: make(map[string]*Logger),
		logger:   base,
	}
}

func (l *Logger) GetLogrusLogger() *logrus.Logger {
	return l.logger
}

func (l *Logger) GetChild(name string) *Logger {
	relative := strings.TrimPrefix(name, l.Name+".")
	abs := strings.TrimPrefix(fmt.Sprintf("%s.%s", l.Name, relative), ".")
	parent := l
	for _, part := range strings.Split(relative, ".") {
		logger, ok := parent.children[part]
		if !ok {
			logger = &Logger{
				Name:     abs,
				parent:   parent,
				absLevel: NilLevel,
				inherit:  true,
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

func (l *Logger) SetLevel(level logrus.Level, inherit bool) error {
	if level == l.absLevel {
		return nil
	}
	if level == NilLevel && l.Name == RootLoggerName {
		return ErrInvalidRootLevel
	}
	l.absLevel = level
	switch level {
	case NilLevel:
		l.logger.Level = l.parent.GetEffectiveLevel()
		l.inherit = true
	default:
		l.logger.Level = level
		l.inherit = inherit
	}
	if inherit {
		l.propagate()
	}
	return nil
}

func (l *Logger) GetEffectiveLevel() logrus.Level {
	if !l.inherit {
		return l.parent.GetEffectiveLevel()
	}
	return l.logger.Level
}

func (l *Logger) GetLevel() logrus.Level {
	return l.absLevel
}

func (l *Logger) propagate() {
	for _, child := range l.children {
		child.inheritLevel(l.GetEffectiveLevel())
	}
}

func (l *Logger) inheritLevel(parentLevel logrus.Level) {
	if l.absLevel == NilLevel {
		l.logger.Level = parentLevel
		l.propagate()
	}
}
