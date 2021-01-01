package snbt

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	testkokens := []token{
		{t: tkArrStart},
		{t: tkMapStart},
		{t: tkStr, v: "firstKey"},
		{t: tkFloat, v: "1.0"},
		{t: tkStr, v: "secondKey"},
		{t: tkFloat, v: "2.0"},
		{t: tkMapEnd},
		{t: tkArrEnd},
	}

	psr := newParser(testkokens)
	outInter, err := psr.parse()
	if err != nil {
		t.Error("failed to parse test tokens: ", err)
	}
	out := outInter.([]interface{})
	if len(out) != 1 {
		t.Errorf("expected a parsed arr of length 1, got %d", len(out))
	}
	firstEl := out[0].(map[string]interface{})
	firstV, ok := firstEl["firstKey"]
	if !ok {
		t.Error("missing 'firstKey' from parsed output")
	}
	if firstV.(float64) != 1.0 {
		t.Errorf("'firstKey' value not properly set, actual=%f, expected=%f", firstV.(float64), 1.0)
	}

	secondV, ok := firstEl["secondKey"]
	if !ok {
		t.Error("missing 'secondKey' from parsed output")
	}
	if secondV.(float64) != 2.0 {
		t.Errorf("'secondKey' value not properly set, actual=%f, expected=%f", secondV.(float64), 2.0)
	}
}

func TestParserWithInvalidTokens(t *testing.T) {
	testInvalidTokens := []token{
		{t: tkArrStart},
		{t: tkArrStart},
		{t: tkFloat, v: "123.1"},
		{t: tkMapEnd},
	}

	psr := newParser(testInvalidTokens)
	_, err := psr.parse()
	if err == nil {
		t.Error("parser failed to error out on invalid tokens")
	}

	if !strings.Contains(err.Error(), "invalid token") {
		t.Error("parser return wrong error: ", err)
	}
	if !strings.Contains(err.Error(), "posn 3") {
		t.Errorf("parser error failed to illegal char at posn 3")
	}
}
