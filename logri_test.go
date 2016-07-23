package logri_test

import (
	"testing"

	"github.com/Sirupsen/logrus/hooks/test"
	. "github.com/iancmcc/logri"
	. "gopkg.in/check.v1"
)

func TestLogri(t *testing.T) { TestingT(t) }

type LogriSuite struct {
	logger *Logger
	hook   *test.Hook
	c      *C
}

type logfunc func(args ...interface{})

var (
	_ = Suite(&LogriSuite{})
)

func (s *LogriSuite) SetUpTest(c *C) {
	logger, hook := test.NewNullLogger()
	s.logger = NewRootLoggerFromLogrus(logger)
	s.hook = hook
	s.c = c
}

func (s *LogriSuite) AssertLogs(f logfunc, m string) {
	defer s.hook.Reset()
	f(m)
	s.c.Assert(s.hook.Entries, HasLen, 1)
	s.c.Assert(s.hook.LastEntry().Message, Equals, m)
}

func (s *LogriSuite) AssertNotLogs(f logfunc, m string) {
	defer s.hook.Reset()
	f(m)
	s.c.Assert(s.hook.Entries, HasLen, 0)
}
