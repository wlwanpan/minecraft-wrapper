package snbt

import (
	"io/ioutil"
	"testing"
)

func TestFullDataGetDecode(t *testing.T) {
	testData, err := ioutil.ReadFile("testdata/data_get_entity")
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

func TestBigDataGetDecode(t *testing.T) {
	testData, err := ioutil.ReadFile("testdata/data_get_entity_big")
	if err != nil {
		t.Errorf("failed to load testdata: %s", err)
		return
	}

	resp := struct {
		Rotation   []float64
		Attributes []struct {
			Name string
			Base float64
		}
		RecipeBook struct {
			Recipes          []string
			IsFurnaceGuiOpen int
			IsGuiOpen        int
		}
	}{}
	if err := Decode(testData, &resp); err != nil {
		t.Errorf("failed to decode testdata: %s", err)
	}
	expectedPos := []float64{-71.11728, 1.9123533}
	for i, v := range expectedPos {
		if resp.Rotation[i] != v {
			t.Errorf("Rotation mismatch value at %d, actual=%f, expected=%f", i, resp.Rotation[i], v)
		}
	}

	expectedAttrs := []struct {
		Name string
		Base float64
	}{
		{"minecraft:generic.attack_damage", 1.0},
		{"minecraft:generic.movement_speed", 0.10000000149011612},
		{"minecraft:generic.attack_speed", 4.0},
	}

	expectedAttrsLen := len(expectedAttrs)
	if len(resp.Attributes) != expectedAttrsLen {
		t.Errorf("Attribute missing items: actual=%d, expected=%d", len(resp.Attributes), expectedAttrsLen)
	}

	for i, attr := range resp.Attributes {
		expectedAttr := expectedAttrs[i]
		if attr.Name != expectedAttr.Name {
			t.Errorf("Attribute 'Name' value mismatch: actual=%s, expected=%s", attr.Name, expectedAttr.Name)
		}
		if attr.Base != expectedAttrs[i].Base {
			t.Errorf("Attribute 'Base' value mismatch: actual=%f, expected=%f", attr.Base, expectedAttr.Base)
		}
	}

	expectedRecipesLen := 482
	if len(resp.RecipeBook.Recipes) != expectedRecipesLen {
		t.Errorf("Attribute missing 'Recipes': actual=%d, expected=%d", len(resp.RecipeBook.Recipes), expectedRecipesLen)
	}

	if resp.RecipeBook.IsFurnaceGuiOpen != 0 {
		t.Errorf("Attribute 'IsFurnaceGuiOpen' should be '0'")
	}

	if resp.RecipeBook.IsGuiOpen != 1 {
		t.Errorf("Attribute 'IsGuiOpen' should be '1'")
	}
}
