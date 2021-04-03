package wrapper

import (
	"bufio"
	"os"
	"testing"

	"github.com/wlwanpan/minecraft-wrapper/events"
)

func testParsedEvents(t *testing.T, evs []events.Event, testfilename string) {
	testfile, err := os.Open(testfilename)
	if err != nil {
		t.Errorf("failed to load test file: %s", err)
	}

	actualEvents := []events.Event{}
	scanner := bufio.NewScanner(testfile)
	for scanner.Scan() {
		ev, t := logParserFunc(scanner.Text(), 0)
		if t == events.TypeNil {
			continue
		}
		actualEvents = append(actualEvents, ev)
	}

	if len(evs) != len(actualEvents) {
		t.Errorf("wrong event count detected: actual=%d, expected=%d", len(actualEvents), len(evs))
		return
	}

	for i, ev := range evs {
		aEv := actualEvents[i]
		if !ev.Is(aEv) {
			t.Errorf("event mismatched at %d: actual=%s, expected=%s", i, aEv.String(), ev.String())
		}
	}
}

func testParsedGameEvents(t *testing.T, gevs []events.GameEvent, testfilename string) {
	testfile, err := os.Open(testfilename)
	if err != nil {
		t.Errorf("failed to load test file: %s", err)
	}

	actualEvents := []events.GameEvent{}
	scanner := bufio.NewScanner(testfile)
	for scanner.Scan() {
		ev, t := logParserFunc(scanner.Text(), 0)
		if t == events.TypeNil {
			continue
		}
		actualEvents = append(actualEvents, ev.(events.GameEvent))
	}

	if len(gevs) != len(actualEvents) {
		t.Logf("actual  : %v", actualEvents)
		t.Logf("expected: %v", gevs)
		t.Errorf("wrong event count detected: actual=%d, expected=%d", len(actualEvents), len(gevs))
		return
	}

	for i, ev := range gevs {
		aEv := actualEvents[i]
		if !ev.Is(aEv) {
			t.Errorf("event mismatched at %d: actual=%s, expected=%s", i, aEv.String(), ev.String())
		}

		for k, v := range ev.Data {
			actualV, ok := aEv.Data[k]
			if !ok {
				t.Errorf("missing data '%s' for event '%s'", v, ev.String())
			}
			if v != actualV {
				t.Errorf("data '%s' mismatch for event '%s': actual=%s, expected=%s ", k, ev.String(), actualV, v)
			}
		}
	}
}

func TestServerStartLog(t *testing.T) {
	evs := []events.Event{
		events.VersionEvent,
		events.StartingEvent,
		events.StartedEvent,
	}
	testParsedEvents(t, evs, "testdata/server_start_log")
}

func TestServerOverloadedLog(t *testing.T) {
	testLagTimes := [][]string{
		{"1", "1"},
		{"2004", "0"},
		{"1000000", "1000000"},
	}

	gevs := []events.GameEvent{}
	for _, times := range testLagTimes {
		gev := events.NewGameEvent(events.ServerOverloaded)
		gev.Data = map[string]string{
			"lag_time": times[0],
			"lag_tick": times[1],
		}
		gevs = append(gevs, gev)
	}

	testParsedGameEvents(t, gevs, "testdata/server_overloaded_log")
}

func TestPlayerBasicLog(t *testing.T) {
	gevs := []events.GameEvent{
		{
			Name: events.PlayerUUID,
			Data: map[string]string{
				"player_name": "player1",
				"player_uuid": "player-1-uuid",
			},
		},
		{
			Name: events.PlayerJoined,
			Data: map[string]string{
				"player_name": "player1",
			},
		},
		{
			Name: events.PlayerUUID,
			Data: map[string]string{
				"player_name": "player2",
				"player_uuid": "player-2-uuid",
			},
		},
		{
			Name: events.PlayerJoined,
			Data: map[string]string{
				"player_name": "player2",
			},
		},
		{
			Name: events.PlayerDied,
			Data: map[string]string{
				"player_name":   "player2",
				"death_by":      "fell",
				"death_details": " from a high place",
			},
		},
		{
			Name: events.PlayerDied,
			Data: map[string]string{
				"player_name":   "player1",
				"death_by":      "was killed by",
				"death_details": " Witch using magic",
			},
		},
		{
			Name: events.PlayerLeft,
			Data: map[string]string{
				"player_name": "player1",
			},
		},
		{
			Name: events.PlayerLeft,
			Data: map[string]string{
				"player_name": "player2",
			},
		},
	}
	testParsedGameEvents(t, gevs, "testdata/player_basic_log")
}
