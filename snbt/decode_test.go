package snbt

import (
	"io/ioutil"
	"testing"
)

func TestFullDataGetDecode(t *testing.T) {
	testData, err := ioutil.ReadFile("testdata/data_get_entity.txt")
	if err != nil {
		t.Errorf("failed to load testdata: %s", err)
		return
	}

	resp := struct {
		HurtByTimestamp int
		Pos             []float64
		Dimension       string
	}{}
	if err := Decode(testData, &resp); err != nil {
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
