package logri_test

import (
	"bytes"

	. "github.com/iancmcc/logri"

	. "gopkg.in/check.v1"
)

var yamlExample = []byte(`
- logger: '*'
  level: info

- logger: 'a.b'
  level: debug


- logger: 'a.b.c.d'
  level: error
  local: true

`)

func getConfig(c *C) LogriConfig {
	r := bytes.NewReader(yamlExample)
	cfg, err := ConfigFromYAML(r)
	c.Assert(err, IsNil)
	return cfg
}

func (s *LogriSuite) TestDeserializeConfig(c *C) {
	cfg := getConfig(c)
	c.Assert(cfg, HasLen, 3)
}

func (s *LogriSuite) TestApplyLevelsFromConfig(chk *C) {
	cfg := getConfig(chk)

	a := s.logger.GetChild("a")         // Should log at info
	b := s.logger.GetChild("a.b")       // debug
	c := s.logger.GetChild("a.b.c")     // debug
	d := s.logger.GetChild("a.b.c.d")   // error
	e := s.logger.GetChild("a.b.c.d.e") // debug

	// Apply the config after the loggers are created
	s.logger.ApplyConfig(cfg)

	s.AssertLogLevel(chk, a, "Info")
	s.AssertLogLevel(chk, b, "Debug")
	s.AssertLogLevel(chk, c, "Debug")
	s.AssertLogLevel(chk, d, "Error")
	s.AssertLogLevel(chk, e, "Debug")
}
