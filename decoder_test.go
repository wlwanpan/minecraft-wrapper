package wrapper

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	testLog := `{Dimension: "minecraft:custom", FakeField: [111, "111"]}`
	lxr := NewLexer([]byte(testLog))
	if err := lxr.Tokenize(); err != nil {
		t.Error("failed to tokenize test str: ", err)
		return
	}
	expectedTokens := []Token{
		{t: TTMapStart},
		{t: TTStr, v: "Dimension"},
		{t: TTStr, v: "minecraft:custom"},
		{t: TTStr, v: "FakeField"},
		{t: TTArrStart},
		{t: TTInt, v: "111"},
		{t: TTStr, v: "111"},
		{t: TTArrEnd},
		{t: TTMapEnd},
	}

	for i, token := range lxr.Tokens() {
		expectedToken := expectedTokens[i]
		if token.t != expectedToken.t {
			t.Errorf("token 't' mismatch at %d, actual=%d, expected=%d", i, token.t, expectedToken.t)
		}
		if token.v != expectedToken.v {
			t.Errorf("token 'v' mismatch at %d, actual=%d, expected=%d", i, token.t, expectedToken.t)
		}
	}
}

func TestParser(t *testing.T) {
	testTokens := []Token{
		{t: TTArrStart},
		{t: TTMapStart},
		{t: TTStr, v: "firstKey"},
		{t: TTFloat, v: "1.0"},
		{t: TTStr, v: "secondKey"},
		{t: TTFloat, v: "2.0"},
		{t: TTMapEnd},
		{t: TTArrEnd},
	}

	psr := NewParser(testTokens)
	outInter, err := psr.Parse()
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
	testInvalidTokens := []Token{
		{t: TTArrStart},
		{t: TTArrStart},
		{t: TTFloat, v: "123.1"},
		{t: TTMapEnd},
	}

	psr := NewParser(testInvalidTokens)
	_, err := psr.Parse()
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

func TestFullDataGetDecode(t *testing.T) {
	testData, err := ioutil.ReadFile("testdata/data_get_entity.txt")
	if err != nil {
		t.Errorf("failed to load testdata: %s", err)
		return
	}

	resp := &DataGetOutput{}
	if err := Decode(testData, resp); err != nil {
		t.Errorf("failed to decode testdata: %s", err)
	}

	expectedHurtByTimestamp := 66261
	if resp.HurtByTimestamp != expectedHurtByTimestamp {
		t.Errorf("HurtByTimestamp no properly set: actual=%d, expected=%d", resp.HurtByTimestamp, expectedHurtByTimestamp)
	}

	expectedPos := []float64{281.30000001192093, 54.0, 367.9814037801891}
	for i, v := range expectedPos {
		if resp.Pos[i] != v {
			t.Errorf("Pos mismatch value at %d, actual=%f, expected=%f", i, resp.Pos[i], v)
		}
	}

	expectedDimension := "minecraft:overworld"
	if resp.Dimension != expectedDimension {
		t.Errorf("Dimension not properly set: actual=%s, expected=%s", resp.Dimension, expectedDimension)
	}
}
