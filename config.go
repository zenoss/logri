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

type LogriConfig struct {
	Loggers map[string]LoggerConfig `yaml:"loggers"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

func ConfigFromYAML(r io.Reader) (cfg LogriConfig, err error) {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	err = yaml.Unmarshal(buf.Bytes(), &cfg)
	return
}
