package logri

import (
	"bytes"
	"errors"
	"io"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	ConfigurationError = errors.New("Unable to parse configuration")
)

// LogriConfig is the configuration for a logri manager
type LogriConfig []LoggerConfig

// LoggerConfig is the configuration for a single logger
type LoggerConfig struct {
	Logger string
	Level  string
	Local  bool
}

func ConfigFromYAML(r io.Reader) (cfg LogriConfig, err error) {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	err = yaml.Unmarshal(buf.Bytes(), &cfg)
	sort.Sort(&cfg)
	return
}

func (c LogriConfig) Len() int      { return len(c) }
func (c LogriConfig) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Sort loggers by place in the hierarchy
func (c LogriConfig) Less(i, j int) bool {
	a, b := c[i].Logger, c[j].Logger
	if a == "*" || a == "" {
		return true
	}
	if b == "*" || b == "" {
		return false
	}
	return strings.Count(a, ".") < strings.Count(b, ".")
}
