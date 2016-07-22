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

// Manager is a manager of multiple loggers.
type Manager struct {
	mu      sync.Mutex
	config  LogriConfig
	loggers map[string]*logrus.Logger
}

// NewManager creates a new log manager.
func NewManager() *Manager {
	return &Manager{loggers: map[string]*logrus.Logger{}}
}

// GetLogger returns the logger with the given name, creating one if necessary.
func (mgr *Manager) GetLogger(name string) (logger *logrus.Logger) {
	var ok bool
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if logger, ok = mgr.loggers[name]; !ok {
		logger = logrus.New()
		mgr.loggers[name] = logger
		mgr.applyConfig()
	}
	return logger
}

// FindLoggers finds loggers with name matching the provided glob pattern.
func (mgr *Manager) FindLoggers(pattern string) (loggers []*logrus.Logger, err error) {
	// Let's not get super picky about this
	if pattern == "*" {
		pattern = "**"
	}
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

func (mgr *Manager) applyConfig() error {
	for pattern, lcfg := range mgr.config.Loggers {
		level, err := logrus.ParseLevel(lcfg.Level)
		if err != nil {
			return err
		}
		loggers, err := mgr.FindLoggers(pattern)
		if err != nil {
			return err
		}
		for _, logger := range loggers {
			logger.Level = level
		}
	}
	return nil

}

// ApplyConfig applies configuration to loggers managed by this manager.
func (mgr *Manager) ApplyConfig(config LogriConfig) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	mgr.config = config
	return mgr.applyConfig()
}

// GetRootLogger returns the root logger.
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

// FindLoggers finds loggers with name matching a glob pattern.
func FindLoggers(pattern string) ([]*logrus.Logger, error) {
	return mgr.FindLoggers(pattern)
}

// ApplyConfig applies configuration to loggers managed by the default manager.
func ApplyConfig(config LogriConfig) error {
	return mgr.ApplyConfig(config)
}
