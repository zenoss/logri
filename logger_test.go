package logri_test

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	. "github.com/iancmcc/logri"
	. "gopkg.in/check.v1"
)

func (s *LogriSuite) TestAThing(c *C) {
	logger := NewRootLogger()
	hook := test.NewLocal(logger.GetLogrusLogger())
	logger.SetLevel(logrus.InfoLevel, true)
	logger.Debug("debug message 1")
	logger.SetLevel(logrus.DebugLevel, true)
	logger.Debug("debug message 2")
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.LastEntry().Message, Equals, "debug message 2")
}

func (s *LogriSuite) TestUnsetLoggerLevel(c *C) {
}

func (s *LogriSuite) TestInheritLevelFromParent(c *C) {
}

func (s *LogriSuite) TestSetRootLoggerToNil(c *C) {
}
