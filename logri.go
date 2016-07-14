package logri

import (
	"errors"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gobwas/glob"
)

var (
	mgr        = NewManager()
	rootLogger = mgr.GetLogger("")
	separator  = '.'

	// ErrInvalidPattern is thrown when the user asks for loggers using an
	// invalid glob pattern
	ErrInvalidPattern = errors.New("Invalid logger name pattern")
)

type Manager struct {
	mu      sync.Mutex
	loggers map[string]*logrus.Logger
}

func NewManager() *Manager {
	return &Manager{loggers: map[string]*logrus.Logger{}}
}

func (mgr *Manager) GetLogger(name string) (logger *logrus.Logger) {
	var ok bool
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if logger, ok = mgr.loggers[name]; !ok {
		logger = logrus.New()
		mgr.loggers[name] = logger
	}
	return logger
}

func (mgr *Manager) FindLoggers(pattern string) (loggers []*logrus.Logger, err error) {
	var compiled glob.Glob
	compiled, err = glob.Compile(pattern, separator)
	if err != nil {
		return
	}
	for k, v := range mgr.loggers {
		if compiled.Match(k) {
			loggers = append(loggers, v)
		}
	}
	return
}

// GetRootLogger returns the root logger
func (mgr *Manager) GetRootLogger() *logrus.Logger {
	return mgr.GetLogger("")
}

// GetLogger returns the logger with the given name, creating one if necessary.
func GetLogger(name string) *logrus.Logger {
	return mgr.GetLogger(name)
}

// GetRootLogger returns the root logger
func GetRootLogger() *logrus.Logger {
	return rootLogger
}

// FindLoggers finds all loggers matching a pattern. You can use *, ** and ?.
func FindLoggers(pattern string) ([]*logrus.Logger, error) {
	return mgr.FindLoggers(pattern)
}
