package logri

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
)

const (
	RootLoggerName              = ""
	MarkerLevel    logrus.Level = 255
	NilLevel       logrus.Level = 254
)

var (
	ErrInvalidRootLevel                    = errors.New("The root logger must have a level")
	RootLogger          logrus.FieldLogger = NewRootLogger()
)

type Logger struct {
	mu         sync.Mutex
	Name       string
	parent     *Logger
	absLevel   logrus.Level
	tmpLevel   logrus.Level
	inherit    bool
	lastConfig LogriConfig
	children   map[string]*Logger
	logger     *logrus.Logger
}

func NewRootLogger() *Logger {
	return NewRootLoggerFromLogrus(logrus.New())
}

func NewRootLoggerFromLogrus(base *logrus.Logger) *Logger {
	return &Logger{
		Name:     RootLoggerName,
		absLevel: logrus.InfoLevel,
		tmpLevel: MarkerLevel,
		inherit:  true,
		children: make(map[string]*Logger),
		logger:   base,
	}
}

func (l *Logger) GetLogrusLogger() *logrus.Logger {
	return l.logger
}

func (l *Logger) GetRoot() *Logger {
	next := l
	for next.parent != nil {
		next = next.parent
	}
	return next
}

func (l *Logger) GetChild(name string) *Logger {
	if name == "" || name == "*" {
		return l.GetRoot()
	}
	relative := strings.TrimPrefix(name, l.Name+".")
	abs := strings.TrimPrefix(fmt.Sprintf("%s.%s", l.Name, relative), ".")
	parent := l
	var changed bool
	for _, part := range strings.Split(relative, ".") {
		logger, ok := parent.children[part]
		if !ok {
			logger = &Logger{
				Name:     abs,
				parent:   parent,
				absLevel: NilLevel,
				tmpLevel: MarkerLevel,
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
			changed = true
		}
		parent = logger
	}
	if changed && l.GetRoot().lastConfig != nil {
		l.ApplyConfig(l.GetRoot().lastConfig)
	}
	return parent
}

func (l *Logger) SetLevel(level logrus.Level, inherit bool) error {
	if err := l.setLevel(level, inherit); err != nil {
		return err
	}
	l.applyTmpLevels()
	return nil
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Out = w
}

func (l *Logger) nilAllLevels() {
	for _, child := range l.children {
		child.absLevel = NilLevel
		child.tmpLevel = MarkerLevel
		child.inherit = true
		child.nilAllLevels()
	}
}

func (l *Logger) setLevel(level logrus.Level, inherit bool) error {
	if level != l.absLevel || l.inherit != inherit {
		if level == NilLevel && l.Name == RootLoggerName {
			return ErrInvalidRootLevel
		}
		l.absLevel = level
		switch level {
		case NilLevel:
			l.tmpLevel = l.parent.GetEffectiveLevel()
			l.inherit = true
		default:
			l.tmpLevel = level
			l.inherit = inherit
		}
	}
	if l.inherit {
		l.propagate()
	}
	return nil
}

func (l *Logger) GetEffectiveLevel() logrus.Level {
	if !l.inherit {
		return l.parent.GetEffectiveLevel()
	}
	if l.tmpLevel != MarkerLevel {
		return l.tmpLevel
	}
	return l.logger.Level
}

func (l *Logger) ApplyConfig(config LogriConfig) error {
	root := l.GetRoot()
	root.nilAllLevels()
	root.lastConfig = config
	// Loggers are already sorted by hierarchy, so we can apply top down safely
	for _, loggerConfig := range config {
		logger := root.GetChild(loggerConfig.Logger)
		level, err := logrus.ParseLevel(loggerConfig.Level)
		if err != nil {
			// TODO: validate before it gets to this point
			return err
		}
		logger.setLevel(level, !loggerConfig.Local)
	}
	root.applyTmpLevels()
	return nil
}

func (l *Logger) propagate() {
	for _, child := range l.children {
		child.inheritLevel(l.GetEffectiveLevel())
	}
}

func (l *Logger) inheritLevel(parentLevel logrus.Level) {
	if l.absLevel == NilLevel {
		l.tmpLevel = parentLevel
		l.propagate()
	}
}

func (l *Logger) applyTmpLevels() {
	if l.tmpLevel != MarkerLevel && l.tmpLevel != l.logger.Level {
		l.logger.Level = l.tmpLevel
	}
	l.tmpLevel = MarkerLevel
	for _, child := range l.children {
		child.applyTmpLevels()
	}
}
