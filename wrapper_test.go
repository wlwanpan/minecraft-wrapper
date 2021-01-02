package wrapper

import (
	"bufio"
	"io"
	"os"
	"testing"
	"time"
)

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

func TestWrapperStart(t *testing.T) {
	c, err := newTestConsole("testdata/server_start_log.txt")
	if err != nil {
		t.Errorf("failed to load test file: %w", err)
		return
	}

	wpr := NewWrapper(c, logParserFunc)
	if wpr.State() != WrapperOffline {
		t.Errorf("wrapper should be 'offline', got %s", wpr.State())
	}

	started, err := wpr.StartAndWait()
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-started:
	case <-time.After(1 * time.Second):
		t.Error("wrapper timeout, failed to start")
	}

	if wpr.State() != WrapperOnline {
		t.Errorf("wrapper should be 'online', got %s", wpr.State())
	}

	expectedDetectedVersion := "1.16.4"
	if wpr.Version != expectedDetectedVersion {
		t.Errorf("wrapper version be %s, got %s", expectedDetectedVersion, wpr.Version)
	}
}
