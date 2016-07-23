package logri_test

import (
	"testing"

	"github.com/Sirupsen/logrus"
	. "github.com/iancmcc/logri"
	. "gopkg.in/check.v1"
)

func TestLogri(t *testing.T) { TestingT(t) }

type LogriSuite struct {
	mgr *Manager
}

var (
	_ = Suite(&LogriSuite{})
)

type loggerListContainsChecker struct {
	*CheckerInfo
}

func (checker *loggerListContainsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	// There are obviously faster ways to do this, but it's extremely unlikely
	// we will have hundreds of thousands of loggers in a test, so whatever,
	// brute force it is
	list := params[0].([]*logrus.Logger)
	therest := params[1].([]*logrus.Logger)
	if len(list) != len(therest) {
		return false, "Results are of different lengths"
	}

	m := make(map[interface{}]int)
	for _, y := range list {
		m[y]++
	}

	for _, x := range therest {
		if m[x] > 0 {
			m[x]--
			continue
		}
		return false, "Logger lists are unequal"
	}

	for _, v := range m {
		if v > 0 {
			return false, "Logger lists are unequal"
		}
	}
	return true, ""
}

var LoggerListContains Checker = &loggerListContainsChecker{
	&CheckerInfo{Name: "LoggerListContains", Params: []string{"obtained", "expected"}},
}

func (s *LogriSuite) SetUpTest(c *C) {
	s.mgr = NewManager()
}

func (s *LogriSuite) TestLoggerCreation(c *C) {
	a := s.mgr.GetLogger("a")
	b := s.mgr.GetLogger("b")

	c.Assert(a, Not(Equals), b)
	c.Assert(a, Equals, s.mgr.GetLogger("a"))
	c.Assert(b, Equals, s.mgr.GetLogger("b"))

	var l *logrus.Logger
	c.Assert(a, FitsTypeOf, l)
	c.Assert(b, FitsTypeOf, l)
}

func (s *LogriSuite) TestFindLoggers(c *C) {
	root := s.mgr.GetRootLogger()
	a := s.mgr.GetLogger("a")
	b := s.mgr.GetLogger("b")
	adb := s.mgr.GetLogger("a.b")
	adbdc := s.mgr.GetLogger("a.b.c")

	loggers, err := s.mgr.FindLoggers("a")
	c.Assert(err, IsNil)
	c.Assert(loggers, LoggerListContains, []*logrus.Logger{a})

	loggers, err = s.mgr.FindLoggers("a.*")
	c.Assert(err, IsNil)
	c.Assert(loggers, LoggerListContains, []*logrus.Logger{adbdc, adb})

	loggers, err = s.mgr.FindLoggers("*")
	c.Assert(err, IsNil)
	c.Assert(loggers, LoggerListContains, []*logrus.Logger{a, b, adbdc, adb, root})

	loggers, err = s.mgr.FindLoggers("")
	c.Assert(err, IsNil)
	c.Assert(loggers, LoggerListContains, []*logrus.Logger{root})
}
