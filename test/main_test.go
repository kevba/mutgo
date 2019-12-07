package main

import "testing"

func TestAdd(t *testing.T) {
	val := add(5, -5)
	if val != 0 {
		t.Errorf("expected 0")
	}
}

func TestIsEqual(t *testing.T) {
	if !isEqual(1, 1) {
		t.Errorf("should be equal")
	}
}
