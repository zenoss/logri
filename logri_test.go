package logri_test

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "github.com/zenoss/logri"
	. "gopkg.in/check.v1"
)

func TestLogri(t *testing.T) { TestingT(t) }

type LogriSuite struct {
	logger *Logger
	hook   *test.Hook
}

type logfunc func(args ...interface{})

var (
	_ = Suite(&LogriSuite{})
)

func (s *LogriSuite) SetUpTest(c *C) {
	logger, hook := test.NewNullLogger()
	s.logger = NewLoggerFromLogrus(logger)
	s.hook = hook
}

func (s *LogriSuite) AssertLogLevel(c *C, ob logrus.FieldLogger, method string) {
	defer s.hook.Reset()

	rob := reflect.ValueOf(ob)
	call := func(m string) {
		meth := rob.MethodByName(m)
		meth.Call([]reflect.Value{reflect.ValueOf(m)})
	}

	// First call the method specified with the method name as the msg
	call(method)

	// Now call one below it to make sure it doesn't log that
	switch method {
	case "Fatal":
		call("Error")
	case "Error":
		call("Warn")
	case "Warn":
		call("Info")
	case "Info":
		call("Debug")
	case "Debug":
	}

	c.Assert(len(s.hook.Entries) > 0, Equals, true, Commentf("No entries were logged at level %s", method))
	c.Assert(len(s.hook.Entries) < 2, Equals, true, Commentf("Entries below %s were logged", method))
	c.Assert(s.hook.LastEntry().Message, Equals, method, Commentf("A level other than %s was logged", method))
}
