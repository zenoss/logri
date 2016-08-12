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
	mu           sync.Mutex
	Name         string
	parent       *Logger
	absLevel     logrus.Level
	tmpLevel     logrus.Level
	inherit      bool
	lastConfig   LogriConfig
	children     map[string]*Logger
	logger       *logrus.Logger
	outputs      []io.Writer
	localOutputs []io.Writer
}

func NewRootLogger() *Logger {
	return NewRootLoggerFromLogrus(logrus.New())
}

func NewRootLoggerFromLogrus(base *logrus.Logger) *Logger {
	return &Logger{
		Name:         RootLoggerName,
		absLevel:     logrus.InfoLevel,
		tmpLevel:     MarkerLevel,
		inherit:      true,
		children:     make(map[string]*Logger),
		logger:       base,
		outputs:      []io.Writer{base.Out},
		localOutputs: []io.Writer{},
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
	parent := l
	var (
		changed  bool
		localabs string = l.Name
	)
	for _, part := range strings.Split(relative, ".") {
		if localabs == "" {
			localabs = part
		} else {
			localabs = fmt.Sprintf("%s.%s", localabs, part)
		}
		logger, ok := parent.children[part]
		if !ok {
			logger = &Logger{
				Name:     localabs,
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
	l.applyTmpState()
	return nil
}

func (l *Logger) addOutput(w io.Writer, inherit bool) {
	if inherit {
		l.outputs = append(l.outputs, w)
	} else {
		l.localOutputs = append(l.localOutputs, w)
	}
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Out = w
}

func (l *Logger) SetOutputs(writers ...io.Writer) {
	l.SetOutput(io.MultiWriter(writers...))
}

func (l *Logger) resetChildren() {
	for _, child := range l.children {
		child.absLevel = NilLevel
		child.tmpLevel = MarkerLevel
		child.inherit = true
		child.outputs = []io.Writer{}
		child.localOutputs = []io.Writer{}
		child.resetChildren()
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
	root.outputs = []io.Writer{}
	root.localOutputs = []io.Writer{}
	root.resetChildren()
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

		for _, outputConfig := range loggerConfig.Out {
			w, err := GetOutputWriter(outputConfig.Type, outputConfig.Options)
			if err != nil {
				return err
			}
			logger.addOutput(w, !outputConfig.Local)
		}
	}
	root.propagate()
	root.applyTmpState()
	return nil
}

func (l *Logger) propagate() {
	for _, child := range l.children {
		child.inheritLevel(l.GetEffectiveLevel())
		child.inheritOutputs(l.getInheritableOutputs())
		child.propagate()
	}
}

func (l *Logger) getInheritableOutputs() []io.Writer {
	var result []io.Writer
	if l.parent != nil {
		for _, out := range l.parent.getInheritableOutputs() {
			result = append(result, out)
		}
	}
	for _, out := range l.outputs {
		result = append(result, out)
	}
	return dedupeWriters(result...)
}

func (l *Logger) inheritOutputs(writers []io.Writer) {
	l.outputs = dedupeWriters(append(l.outputs, writers...)...)
}

func (l *Logger) inheritLevel(parentLevel logrus.Level) {
	if l.absLevel == NilLevel {
		l.tmpLevel = parentLevel
	}
}

func (l *Logger) applyTmpState() {
	if l.tmpLevel != MarkerLevel && l.tmpLevel != l.logger.Level {
		l.logger.Level = l.tmpLevel
	}
	l.tmpLevel = MarkerLevel
	allwriters := append(l.outputs, l.localOutputs...)
	l.SetOutputs(dedupeWriters(allwriters...)...)
	for _, child := range l.children {
		child.applyTmpState()
	}
}

func dedupeWriters(writers ...io.Writer) []io.Writer {
	var val struct{}
	m := map[io.Writer]struct{}{}
	for _, writer := range writers {
		m[writer] = val
	}
	var result []io.Writer
	for writer := range m {
		result = append(result, writer)
	}
	return result
}
