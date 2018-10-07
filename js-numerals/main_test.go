package main

import "testing"

func TestOneDigit(t *testing.T) {
	if result := number_to_string(5); "five" != result {
		t.Error("Expected five, got ", result)
	}
}
