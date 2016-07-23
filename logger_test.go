package logri_test

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	. "github.com/iancmcc/logri"
	. "gopkg.in/check.v1"
)

func (s *LogriSuite) TestAThing(c *C) {
	logger := Logger{}
	hook := test.NewLocal(logger)
	logger.SetLevel(logrus.InfoLevel)
	c.Assert(logger.Level, Equals, logrus.InfoLevel)
}

func (s *LogriSuite) TestUnsetLoggerLevel(c *C) {
}

func (s *LogriSuite) TestInheritLevelFromParent(c *C) {
}

func (s *LogriSuite) TestSetRootLoggerToNil(c *C) {
	logger := Logger{}
	logger.SetLevel(logrus.InfoLevel)
	c.Assert(logger.GetLevel(), Equals, logrus.InfoLevel)
}
