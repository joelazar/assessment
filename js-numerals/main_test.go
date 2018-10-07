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

func TestThreeDigit(t *testing.T) {
	if result := number_to_string(101); "one hundred and one" != result {
		t.Error("Expected one hundred and one, got ", result)
	}
}

func TestHigherThreeDigit(t *testing.T) {
	if result := number_to_string(982); "nine hundred and eighty-two" != result {
		t.Error("Expected nine hundred and eighty-two, got ", result)
	}
}
