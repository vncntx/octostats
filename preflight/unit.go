package preflight

import (
	"io/ioutil"
	"os"
	"testing"
)

// UnitTest provides utilities for unit testing
type UnitTest struct {
	*testing.T
}

// Action is a function with no arguments
type Action func()

// FileConsumer is a function that consumes a file
type FileConsumer func(file *os.File)

// Unit returns a new unit test
func Unit(t *testing.T) *UnitTest {
	return &UnitTest{t}
}

// Expect returns a new value-based expectation
func (unit *UnitTest) Expect(actual interface{}) Expectation {
	return ExpectValue(unit.T, actual)
}

// ExpectFile returns expectations based on file contents
func (unit *UnitTest) ExpectFile(file *os.File) Expectation {
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		unit.Errorf("could not read from file '%s'", file.Name())
	}

	return unit.Expect(string(contents))
}

// ExpectOutput returns a new output file-based expectation
func (unit *UnitTest) ExpectOutput(consumer FileConsumer) Expectation {
	read, write := unit.createStream()

	// invoke the consumer
	consumer(write)
	unit.closeStream(write)

	return unit.ExpectFile(read)
}

// ExpectExitCode overrides the scaffolding osExit function
func (unit *UnitTest) ExpectExitCode(act Action) Expectation {
	exitCode := 0
	Scaffold.OSExit = func(code int) {
		exitCode = code
	}
	act()
	Restore()

	return unit.Expect(exitCode)
}

func (unit *UnitTest) createStream() (readable *os.File, writable *os.File) {
	readable, writable, err := os.Pipe()
	if err != nil {
		unit.Errorf("failed to create stream: %s", err)
	}

	return
}

func (unit *UnitTest) closeStream(stream *os.File) {
	if err := stream.Close(); err != nil {
		unit.Errorf("failed to close stream %s: %s", stream.Name(), err)
	}
}
