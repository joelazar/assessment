package main

import "testing"

func TestOneDigit(t *testing.T) {
	if result := number_to_string(5); "five" != result {
		t.Error("Expected five, got ", result)
	}
}

func TestTwoDigit(t *testing.T) {
	if result := number_to_string(13); "thirteen" != result {
		t.Error("Expected thirteen, got ", result)
	}
}

func TestHigherTwoDigit(t *testing.T) {
	if result := number_to_string(97); "ninety-seven" != result {
		t.Error("Expected ninety-seven, got ", result)
	}
}
