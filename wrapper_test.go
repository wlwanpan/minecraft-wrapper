package wrapper

import (
	"testing"
	"time"
)

func TestWrapperStart(t *testing.T) {
	c, err := newTestConsole("testdata/server_start_log")
	if err != nil {
		t.Errorf("failed to load test file: %w", err)
		return
	}

	wpr := NewWrapper(c, logParserFunc)
	if wpr.State() != WrapperOffline {
		t.Errorf("wrapper should be 'offline', got %s", wpr.State())
	}

	if err := wpr.Start(); err != nil {
		t.Error(err)
		return
	}
	select {
	case <-wpr.Loaded():
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

func TestWrapperOffline(t *testing.T) {
	c, err := newTestConsole("testdata/server_start_log")
	if err != nil {
		t.Errorf("failed to load test file: %w", err)
		return
	}

	wpr := NewWrapper(c, logParserFunc)
	if wpr.State() != WrapperOffline {
		t.Errorf("wrapper should be 'offline', got %s", wpr.State())
	}

	// test simple entry command (no output expected).
	if err := wpr.Ban("player-1", "reason-1"); err == nil {
		t.Error("wrapper.Ban should error when 'offline'")
	}

	// test single entry command (with expected output).
	_, err = wpr.DataGet("entity", "player-1")
	if err == nil {
		t.Error("wrapper.DataGet should error when 'offline'")
	}

	// test list entry command (with output parsing multiple log lines).
	_, err = wpr.BanList(BanPlayers)
	if err == nil {
		t.Error("wrapper.BanList should error when 'offline'")
	}
}
