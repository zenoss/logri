package logri_test

import (
	"github.com/Sirupsen/logrus"
	. "github.com/iancmcc/logri"
	. "gopkg.in/check.v1"
)

func (s *LogriSuite) TestDefault(c *C) {
	cfg := LogriConfig{Loggers: map[string]LoggerConfig{
		"*":            LoggerConfig{Level: "fatal"},
		"something.**": LoggerConfig{Level: "debug"},
	}}
	logger1 := s.mgr.GetLogger("hello")
	logger2 := s.mgr.GetLogger("hello.world")
	logger3 := s.mgr.GetLogger("something.else")

	s.mgr.ApplyConfig(cfg)

	c.Assert(logger1.Level, Equals, logrus.FatalLevel)
	c.Assert(logger2.Level, Equals, logrus.FatalLevel)
	c.Assert(logger3.Level, Equals, logrus.DebugLevel)
}

func (s *LogriSuite) TestConfigOrder(c *C) {
	cfg := LogriConfig{Loggers: map[string]LoggerConfig{
		"*.world": LoggerConfig{Level: "debug"},
		"*":       LoggerConfig{Level: "fatal"},
		"hello":   LoggerConfig{Level: "error"},
	}}
	logger1 := s.mgr.GetLogger("hello")
	logger2 := s.mgr.GetLogger("hello.world")
	logger3 := s.mgr.GetLogger("something.else")

	s.mgr.ApplyConfig(cfg)

	c.Assert(logger1.Level, Equals, logrus.ErrorLevel)
	c.Assert(logger2.Level, Equals, logrus.FatalLevel)
	c.Assert(logger3.Level, Equals, logrus.FatalLevel)
}

func (s *LogriSuite) TestConfigAppliesToNewLoggers(c *C) {
	cfg := LogriConfig{Loggers: map[string]LoggerConfig{
		"*":       LoggerConfig{Level: "debug"},
		"hello.*": LoggerConfig{Level: "error"},
	}}

	logger1 := s.mgr.GetLogger("whatever")
	s.mgr.ApplyConfig(cfg)
	c.Assert(logger1.Level, Equals, logrus.DebugLevel)

	logger2 := s.mgr.GetLogger("hello.world")
	c.Assert(logger2.Level, Equals, logrus.ErrorLevel)
}
