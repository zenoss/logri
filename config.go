package logri

import (
	"bytes"
	"errors"
	"io"

	"gopkg.in/yaml.v2"
)

var (
	ConfigurationError = errors.New("Unable to parse configuration")
)

// LogriConfig is the configuration for a logri manager
type LogriConfig struct {
	Loggers map[string]LoggerConfig `yaml:"loggers"`
}

// LoggerConfig is the configuration for a single logger
type LoggerConfig struct {
	Level string `yaml:"level"`
}

func ConfigFromYAML(r io.Reader) (cfg LogriConfig, err error) {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	err = yaml.Unmarshal(buf.Bytes(), &cfg)
	return
}
