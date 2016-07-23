package logri_test

import (
	"github.com/Sirupsen/logrus"
	. "github.com/iancmcc/logri"
	. "gopkg.in/check.v1"
)

func (s *LogriSuite) TestAThing(c *C) {
	logger := Logger{}
	logger.SetLevel(logrus.InfoLevel)
	c.Assert(logger.GetLevel(), Equals, logrus.InfoLevel)
}
