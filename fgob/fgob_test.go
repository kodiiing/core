package fgob_test

import (
	"kodiiing/fgob"
	"testing"
)

func TestIntegration(t *testing.T) {
	type Point struct {
		X int64
		Y int64
	}

	type Store struct {
		ID          uint32
		Name        string
		Open        bool
		Coordinates Point
		Menu        []string
	}

	store := Store{
		ID:   1,
		Name: "Test",
		Open: true,
		Coordinates: Point{
			X: 1,
			Y: 2,
		},
		Menu: []string{"a", "b", "c"},
	}

	data, err := fgob.Marshal(store)
	if err != nil {
		t.Errorf("error marshaling: %v", err)
	}

	if len(data) <= 0 {
		t.Errorf("data is empty")
	}

	var store2 Store
	err = fgob.Unmarshal(data, &store2)
	if err != nil {
		t.Errorf("error unmarshaling: %v", err)
	}

	if store2.ID != store.ID {
		t.Errorf("id is not equal")
	}

	if store2.Name != store.Name {
		t.Errorf("name is not equal")
	}

	if store2.Open != store.Open {
		t.Errorf("open is not equal")
	}

	if store2.Coordinates.X != store.Coordinates.X {
		t.Errorf("x is not equal")
	}

	if store2.Coordinates.Y != store.Coordinates.Y {
		t.Errorf("y is not equal")
	}

	if len(store2.Menu) != len(store.Menu) {
		t.Errorf("menu is not equal")
	}
}

func TestMarshalError(t *testing.T) {
	_, err := fgob.Marshal(nil)
	if err == nil {
		t.Errorf("error expected")
	}
}
