package logri_test

import (
	"bytes"

	"github.com/sirupsen/logrus"
	. "github.com/zenoss/logri"
	. "gopkg.in/check.v1"
)

var (
	nilLevel logrus.Level = 254
)

func (s *LogriSuite) TestSetLoggerLevel(c *C) {
	s.logger.SetLevel(logrus.InfoLevel, true)
	s.AssertLogLevel(c, s.logger, "Info")

	s.logger.SetLevel(logrus.DebugLevel, true)
	s.AssertLogLevel(c, s.logger, "Debug")
}

func (s *LogriSuite) TestUnsetLoggerLevel(c *C) {
	err := s.logger.SetLevel(nilLevel, true)
	c.Assert(err.Error(), Equals, ErrInvalidRootLevel.Error())
	alogger := s.logger.GetChild("a")
	alogger.SetLevel(logrus.ErrorLevel, true)
	s.AssertLogLevel(c, alogger, "Error")
	err = alogger.SetLevel(nilLevel, true)
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
	d.SetLevel(nilLevel, true)
	s.AssertLogLevel(chk, d, "Debug")
	s.AssertLogLevel(chk, e, "Debug")

	// Set c's level to nilLevel, which it already is. Shouldn't affect anything
	c.SetLevel(nilLevel, true)
	s.AssertLogLevel(chk, c, "Debug")
	s.AssertLogLevel(chk, d, "Debug")

	// Now set c's level to something else and back to Nil. Still should
	// inherit from root
	c.SetLevel(logrus.FatalLevel, true)
	c.SetLevel(nilLevel, true)
	s.AssertLogLevel(chk, c, "Debug")
	s.AssertLogLevel(chk, d, "Debug")
	s.AssertLogLevel(chk, e, "Debug")
}

func (s *LogriSuite) TestSetOutput(c *C) {
	var w bytes.Buffer
	s.logger.Info("message")
	c.Assert(w.Len(), Equals, 0)
	s.logger.SetOutput(&w)
	s.logger.Info("message 2")
	c.Assert(w.Len() > 0, Equals, true)
}

func (s *LogriSuite) TestSetOutputs(c *C) {
	var w1, w2 bytes.Buffer

	s.logger.Info("message")
	c.Assert(w1.Len(), Equals, 0)
	c.Assert(w2.Len(), Equals, 0)

	w1.Reset()
	w2.Reset()

	s.logger.SetOutputs(&w1)
	s.logger.Info("message 2")
	c.Assert(w1.Len() > 0, Equals, true)
	c.Assert(w2.Len() > 0, Equals, false)

	w1.Reset()
	w2.Reset()

	s.logger.SetOutputs(&w1, &w2)
	s.logger.Info("message 3")
	c.Assert(w1.Len() > 0, Equals, true)
	c.Assert(w2.Len() > 0, Equals, true)
}
