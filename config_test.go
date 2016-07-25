package logri_test

import (
	"bytes"

	. "github.com/iancmcc/logri"

	. "gopkg.in/check.v1"
)

var inOrder = []byte(`
- logger: '*'
  level: info

- logger: 'a.b'
  level: debug


- logger: 'a.b.c.d'
  level: error
  local: true
`)

var inOrderMinusOne = []byte(`
- logger: '*'
  level: info

- logger: 'a.b'
  level: debug
`)

var inOrderAllInherited = []byte(`
- logger: '*'
  level: info

- logger: 'a.b'
  level: debug


- logger: 'a.b.c.d'
  level: error
`)

var outOfOrder = []byte(`
- logger: 'a.b.c.d'
  level: error
  local: true

- logger: 'a.b'
  level: debug

- logger: '*'
  level: info
`)

var reset = []byte(`
- logger: '*'
  level: info
`)

func getConfig(c *C, yaml []byte) LogriConfig {
	r := bytes.NewReader(yaml)
	cfg, err := ConfigFromYAML(r)
	c.Assert(err, IsNil)
	return cfg
}

func (s *LogriSuite) TestDeserializeConfig(c *C) {
	cfg := getConfig(c, inOrder)
	c.Assert(cfg, HasLen, 3)
}

func (s *LogriSuite) TestApplyLevelsFromConfig(chk *C) {
	cfg := getConfig(chk, inOrder)

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

func (s *LogriSuite) TestApplyLevelsOutOfOrder(chk *C) {
	cfg := getConfig(chk, outOfOrder)

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

func (s *LogriSuite) TestUnsetLocalLevel(chk *C) {
	cfg1 := getConfig(chk, inOrder)
	cfg2 := getConfig(chk, inOrderAllInherited)

	a := s.logger.GetChild("a")         // Should log at info
	b := s.logger.GetChild("a.b")       // debug
	c := s.logger.GetChild("a.b.c")     // debug
	d := s.logger.GetChild("a.b.c.d")   // error
	e := s.logger.GetChild("a.b.c.d.e") // debug

	// Apply the config after the loggers are created
	s.logger.ApplyConfig(cfg1)

	s.AssertLogLevel(chk, a, "Info")
	s.AssertLogLevel(chk, b, "Debug")
	s.AssertLogLevel(chk, c, "Debug")
	s.AssertLogLevel(chk, d, "Error")
	s.AssertLogLevel(chk, e, "Debug")

	// Apply the second config, which should unset the local property on d,
	// causing its level to be inherited by e
	s.logger.ApplyConfig(cfg2)

	s.AssertLogLevel(chk, a, "Info")
	s.AssertLogLevel(chk, b, "Debug")
	s.AssertLogLevel(chk, c, "Debug")
	s.AssertLogLevel(chk, d, "Error")
	s.AssertLogLevel(chk, e, "Error")

}

func (s *LogriSuite) TestAddLoggerConfig(chk *C) {
	cfg1 := getConfig(chk, inOrderMinusOne)
	cfg2 := getConfig(chk, inOrder)

	a := s.logger.GetChild("a")         // Should log at info
	b := s.logger.GetChild("a.b")       // debug
	c := s.logger.GetChild("a.b.c")     // debug
	d := s.logger.GetChild("a.b.c.d")   // error
	e := s.logger.GetChild("a.b.c.d.e") // debug

	// Apply the config after the loggers are created
	s.logger.ApplyConfig(cfg1)

	s.AssertLogLevel(chk, a, "Info")
	s.AssertLogLevel(chk, b, "Debug")
	s.AssertLogLevel(chk, c, "Debug")
	s.AssertLogLevel(chk, d, "Debug")
	s.AssertLogLevel(chk, e, "Debug")

	// Apply the second config, which should unset the local property on d,
	// causing its level to be inherited by e
	s.logger.ApplyConfig(cfg2)

	s.AssertLogLevel(chk, a, "Info")
	s.AssertLogLevel(chk, b, "Debug")
	s.AssertLogLevel(chk, c, "Debug")
	s.AssertLogLevel(chk, d, "Error")
	s.AssertLogLevel(chk, e, "Debug")

}

func (s *LogriSuite) TestDeleteLoggerConfig(chk *C) {
	cfg1 := getConfig(chk, inOrder)
	cfg2 := getConfig(chk, reset)

	a := s.logger.GetChild("a")         // Should log at info
	b := s.logger.GetChild("a.b")       // debug
	c := s.logger.GetChild("a.b.c")     // debug
	d := s.logger.GetChild("a.b.c.d")   // error
	e := s.logger.GetChild("a.b.c.d.e") // debug

	// Apply the config after the loggers are created
	s.logger.ApplyConfig(cfg1)

	s.AssertLogLevel(chk, a, "Info")
	s.AssertLogLevel(chk, b, "Debug")
	s.AssertLogLevel(chk, c, "Debug")
	s.AssertLogLevel(chk, d, "Error")
	s.AssertLogLevel(chk, e, "Debug")

	// Apply the second config, which should unset the local property on d,
	// causing its level to be inherited by e
	s.logger.ApplyConfig(cfg2)

	s.AssertLogLevel(chk, a, "Info")
	s.AssertLogLevel(chk, b, "Info")
	s.AssertLogLevel(chk, c, "Info")
	s.AssertLogLevel(chk, d, "Info")
	s.AssertLogLevel(chk, e, "Info")

}
