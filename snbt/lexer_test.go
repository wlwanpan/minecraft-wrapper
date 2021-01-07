package snbt

import (
	"testing"
)

func TestBasicLexer(t *testing.T) {
	testLog := `{Dimension: "minecraft:custom", FakeField: [111, "111"], Rotation: [-1000, 1000]}`
	lxr := newLexer([]byte(testLog))
	if err := lxr.tokenize(); err != nil {
		t.Error("failed to tokenize test str: ", err)
		return
	}
	expectedTokens := []token{
		{t: tkMapStart},
		{t: tkStr, v: "Dimension"},
		{t: tkStr, v: "minecraft:custom"},
		{t: tkStr, v: "FakeField"},
		{t: tkArrStart},
		{t: tkInt, v: "111"},
		{t: tkStr, v: "111"},
		{t: tkArrEnd},
		{t: tkStr, v: "Rotation"},
		{t: tkArrStart},
		{t: tkInt, v: "-1000"},
		{t: tkInt, v: "1000"},
		{t: tkArrEnd},
		{t: tkMapEnd},
	}

	for i, token := range lxr.tokens() {
		expectedToken := expectedTokens[i]
		if token.t != expectedToken.t {
			t.Errorf("token 't' mismatch at %d, actual=%d, expected=%d", i, token.t, expectedToken.t)
		}
		if token.v != expectedToken.v {
			t.Errorf("token 'v' mismatch at %d, actual=%d, expected=%d", i, token.t, expectedToken.t)
		}
	}
}
