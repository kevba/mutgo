package main

import "testing"

func TestAdd(t *testing.T) {
	val := add(2, 2)
	if val != 4 {
		t.Errorf("expected 0")
	}
}

func TestIsEqual(t *testing.T) {
	if !isEqual(1, 1) {
		t.Errorf("should be equal")
	}
}
