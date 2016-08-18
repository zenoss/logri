package logri_test

import (
	"bytes"

	. "github.com/zenoss/logri"

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

var simplebuffer = []byte(`
- logger: '*'
  level: info
  out:
  - type: test
    options:
        name: test1
`)

var complexbuffers = []byte(`
- logger: '*'
  level: info
  out:
  - type: test
    options:
        name: root
  - type: test
    options:
        name: root2
- logger: a
  level: debug
  out:
  - type: test
    options:
      name: abuf
  - type: test
    local: true
    options:
      name: abuflocal
- logger: a.b
  level: warn
`)

func getConfig(c *C, yaml []byte) LogriConfig {
	r := bytes.NewReader(yaml)
	cfg, err := ConfigFromYAML(r)
	c.Assert(err, IsNil)
	return cfg
}

func getOutputBufferNamed(name string) *bytes.Buffer {
	buffer, _ := GetOutputWriter(TestOutput, map[string]string{"name": name})
	return buffer.(*bytes.Buffer)
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

func (s *LogriSuite) TestLoggersUseExistingConfig(c *C) {
	cfg := getConfig(c, inOrder)
	a := s.logger.GetChild("a")   // Should log at info
	b := s.logger.GetChild("a.b") // debug

	// Apply the config after the loggers are created
	s.logger.ApplyConfig(cfg)

	f := s.logger.GetChild("a.b.c.d.e.f") // error

	s.AssertLogLevel(c, a, "Info")
	s.AssertLogLevel(c, b, "Debug")
	s.AssertLogLevel(c, f, "Debug")

}

func (s *LogriSuite) TestSimpleOutputConfig(c *C) {
	cfg := getConfig(c, simplebuffer)

	s.logger.ApplyConfig(cfg)

	buf := getOutputBufferNamed("test1")
	defer buf.Reset()

	c.Assert(buf.Len(), Equals, 0)
	s.logger.Info("HI")
	c.Assert(buf.Len(), Not(Equals), 0)
}

func (s *LogriSuite) TestSimpleOutputConfigInherited(c *C) {
	cfg := getConfig(c, simplebuffer)
	a := s.logger.GetChild("a")
	s.logger.ApplyConfig(cfg)

	buf := getOutputBufferNamed("test1")
	defer buf.Reset()

	c.Assert(buf.Len(), Equals, 0)
	a.Info("HI")
	c.Assert(buf.Len(), Not(Equals), 0)
}

func (s *LogriSuite) TestComplexOutputBuffer(c *C) {
	cfg := getConfig(c, complexbuffers)
	ab := s.logger.GetChild("a.b")

	s.logger.ApplyConfig(cfg)

	rootbuf := getOutputBufferNamed("root")
	root2buf := getOutputBufferNamed("root2")
	abuf := getOutputBufferNamed("abuf")
	abuflocal := getOutputBufferNamed("abuflocal")

	reset := func() {
		rootbuf.Reset()
		root2buf.Reset()
		abuf.Reset()
		abuflocal.Reset()
	}

	reset()
	defer reset()

	ab.Error("TEST") // Should write to rootbuf, root2buf, abuf

	c.Assert(rootbuf.Len(), Not(Equals), 0)
	c.Assert(root2buf.Len(), Not(Equals), 0)
	c.Assert(abuf.Len(), Not(Equals), 0)
	c.Assert(abuflocal.Len(), Equals, 0)

	reset()

	a := s.logger.GetChild("a")
	a.Debug("TEST") // Should write to all buffers

	c.Assert(abuf.Len(), Not(Equals), 0)
	c.Assert(abuflocal.Len(), Not(Equals), 0)
	c.Assert(rootbuf.Len(), Not(Equals), 0)
	c.Assert(root2buf.Len(), Not(Equals), 0)

}
