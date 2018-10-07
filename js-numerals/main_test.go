package main

import "testing"

func testNumberConverting(t *testing.T, number int, expected_string string) {
	if result := number_to_string(number); expected_string != result {
		t.Errorf("Expected %s, but got %s", expected_string, result)
	}
}

func TestOneDigit(t *testing.T) {
	testNumberConverting(t, 5, "five")
}

func TestTwoDigit(t *testing.T) {
	testNumberConverting(t, 13, "thirteen")
}

func TestHigherTwoDigit(t *testing.T) {
	testNumberConverting(t, 97, "ninety-seven")
}

func TestThreeDigit(t *testing.T) {
	testNumberConverting(t, 101, "one hundred and one")
}

func TestHigherThreeDigit(t *testing.T) {
	testNumberConverting(t, 982, "nine hundred and eighty-two")
}

func TestFourDigit(t *testing.T) {
	testNumberConverting(t, 1001, "one thousand and one")
}

func TestNineDigit(t *testing.T) {
	testNumberConverting(t, 123456789, "one hundred and twenty-three million and four hundred and fifty-six thousand and seven hundred and eighty-nine")
}

func TestMinusNumber(t *testing.T) {
	testNumberConverting(t, -123, "minus one hundred and twenty-three")
}
