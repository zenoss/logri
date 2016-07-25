package logri_test

import (
	"github.com/Sirupsen/logrus"
	. "github.com/iancmcc/logri"
	. "gopkg.in/check.v1"
)

func (s *LogriSuite) TestSetLoggerLevel(c *C) {
	s.logger.SetLevel(logrus.InfoLevel, true)
	s.AssertLogLevel(c, s.logger, "Info")

	s.logger.SetLevel(logrus.DebugLevel, true)
	s.AssertLogLevel(c, s.logger, "Debug")
}

func (s *LogriSuite) TestUnsetLoggerLevel(c *C) {
	err := s.logger.SetLevel(NilLevel, true)
	c.Assert(err.Error(), Equals, ErrInvalidRootLevel.Error())
	alogger := s.logger.GetChild("a")
	alogger.SetLevel(logrus.ErrorLevel, true)
	s.AssertLogLevel(c, alogger, "Error")
	err = alogger.SetLevel(NilLevel, true)
	c.Assert(err, IsNil)
	s.AssertLogLevel(c, alogger, "Info")
}

func (s *LogriSuite) TestGetChildLogger(c *C) {
	alogger := s.logger.GetChild("a")
	c.Assert(alogger.Name, Equals, "a")

	blogger := alogger.GetChild("b")
	c.Assert(blogger.Name, Equals, "a.b")

	blogger2 := s.logger.GetChild("a.b")
	c.Assert(blogger2, Equals, blogger)

	clogger := s.logger.GetChild("a.b.c")
	c.Assert(clogger.Name, Equals, "a.b.c")

	clogger2 := alogger.GetChild("a.b.c")
	c.Assert(clogger2, Equals, clogger)

	clogger3 := blogger.GetChild("a.b.c")
	c.Assert(clogger3, Equals, clogger)

	clogger4 := blogger.GetChild("c")
	c.Assert(clogger4, Equals, clogger)

	clogger5 := alogger.GetChild("b.c")
	c.Assert(clogger5, Equals, clogger)

	clogger6 := blogger.GetChild("d.b.c")
	c.Assert(clogger6, Not(Equals), clogger)
}

func (s *LogriSuite) TestInheritLevelFromParent(chk *C) {
	a := s.logger.GetChild("a")
	b := s.logger.GetChild("a.b")
	c := s.logger.GetChild("a.b.c")
	d := s.logger.GetChild("a.b.c.d")
	e := s.logger.GetChild("a.b.c.d.e")

	s.logger.SetLevel(logrus.DebugLevel, true)
	b.SetLevel(logrus.ErrorLevel, false) // Don't propagate
	d.SetLevel(logrus.InfoLevel, true)

	s.AssertLogLevel(chk, a, "Debug")
	s.AssertLogLevel(chk, b, "Error")
	s.AssertLogLevel(chk, c, "Debug")
	s.AssertLogLevel(chk, d, "Info")
	s.AssertLogLevel(chk, e, "Info")

	// Unset d's level. Now d and e should inherit from the root, since b is a
	// non-propagate level
	d.SetLevel(NilLevel, true)
	s.AssertLogLevel(chk, d, "Debug")
	s.AssertLogLevel(chk, e, "Debug")

	// Set c's level to NilLevel, which it already is. Shouldn't affect anything
	c.SetLevel(NilLevel, true)
	s.AssertLogLevel(chk, c, "Debug")
	s.AssertLogLevel(chk, d, "Debug")

	// Now set c's level to something else and back to Nil. Still should
	// inherit from root
	c.SetLevel(logrus.FatalLevel, true)
	c.SetLevel(NilLevel, true)
	s.AssertLogLevel(chk, c, "Debug")
	s.AssertLogLevel(chk, d, "Debug")
	s.AssertLogLevel(chk, e, "Debug")
}
