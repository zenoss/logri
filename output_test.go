package logri_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/zenoss/logri"

	. "gopkg.in/check.v1"
)

var emptyOptions map[string]string

func (s *LogriSuite) TestStdGetOutputWriters(c *C) {
	w, err := GetOutputWriter(StdoutOutput, emptyOptions)
	c.Assert(err, IsNil)
	c.Assert(w, Equals, os.Stdout)

	w, err = GetOutputWriter(StderrOutput, emptyOptions)
	c.Assert(err, IsNil)
	c.Assert(w, Equals, os.Stderr)
}

func (s *LogriSuite) TestGetFileOutputWriters(c *C) {
	dir := c.MkDir()

	file1 := filepath.Join(dir, "file1")
	file2 := filepath.Join(dir, "file2")

	// Validate that failing to pass a file option is an error
	w1, err := GetOutputWriter(FileOutput, emptyOptions)
	c.Assert(err, Equals, ErrInvalidOutputOptions)

	// Get a proper writer
	w1, err = GetOutputWriter(FileOutput, map[string]string{"file": file1})
	c.Assert(err, IsNil)

	// Verify that we actually have a writer to the file we specified, and that
	// it was properly created
	data := []byte("these are some data")
	w1.Write(data)
	read, err := ioutil.ReadFile(file1)
	c.Assert(err, IsNil)
	c.Assert(string(data), Equals, string(read))

	// Make sure that a different file gives a different writer
	w2, err := GetOutputWriter(FileOutput, map[string]string{"file": file2})
	data2 := []byte("it other data now")
	w2.Write(data2)
	read2, err := ioutil.ReadFile(file2)
	c.Assert(err, IsNil)
	c.Assert(string(data2), Equals, string(read2))
}
