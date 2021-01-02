package wrapper

import (
	"bufio"
	"io"
	"os"
)

// testConsole provide a test console implementation of the interface Console,
// that reads from a test log file instead of a running java stdout. This is
// mainly due used for unit tests in test files.
type testConsole struct {
	scnr *bufio.Scanner
}

func (tc *testConsole) Start() error {
	return nil
}

func (tc *testConsole) Kill() error {
	return nil
}

func (tc *testConsole) WriteCmd(c string) error {
	return nil
}

func (tc *testConsole) ReadLine() (string, error) {
	if tc.scnr.Scan() {
		return tc.scnr.Text(), nil
	}
	return "", io.EOF
}

func newTestConsole(filename string) (*testConsole, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &testConsole{
		scnr: bufio.NewScanner(file),
	}, nil
}
