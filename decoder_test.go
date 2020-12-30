package wrapper

import (
	"io/ioutil"
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
