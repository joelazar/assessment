package main

import (
	"fmt"
	"os"
	"strconv"
)

var (
	smallnumbers = map[int]string{
		0:  "zero",
		1:  "one",
		2:  "two",
		3:  "three",
		4:  "four",
		5:  "five",
		6:  "six",
		7:  "seven",
		8:  "eight",
		9:  "nine",
		10: "ten",
		11: "eleven",
		12: "twelve",
		13: "thirteen",
		14: "fourteen",
		15: "fifteen",
		16: "sixteen",
		17: "seventeen",
		18: "eightteen",
		19: "nineteen",
	}

	tens = map[int]string{
		20: "twenty",
		30: "thirty",
		40: "forty",
		50: "fifty",
		60: "sixty",
		70: "seventy",
		80: "eighty",
		90: "ninety",
	}

	scalenumbers = map[int]string{
		100:        "hundred",
		1000:       "thousand",
		1000000:    "million",
		1000000000: "billion",
	}
)

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}

func scalenumber_to_string(number int, scale int) string {
	quotient, remainder := divmod(number, scale)
	result := fmt.Sprintf("%s %s", number_to_string(quotient), scalenumbers[scale])
	if remainder != 0 {
		result += " and "
		number -= quotient * scale
		result += number_to_string(number)
	}
	return result
}

func number_to_string(number int) string {
	result := ""
	if number > 1000000000 {
		result += scalenumber_to_string(number, 1000000000)
	} else if number > 1000000 {
		result += scalenumber_to_string(number, 1000000)
	} else if number > 1000 {
		result += scalenumber_to_string(number, 1000)
	} else if number > 100 {
		result += scalenumber_to_string(number, 100)
	} else if number > 19 {
		_, remainder := divmod(number, 10)
		if remainder != 0 {
			result += fmt.Sprintf("%s-%s", tens[number-remainder], smallnumbers[remainder])
		} else {
			result += fmt.Sprintf("%s", tens[number-remainder])
		}
	} else if number < 0 {
		result += "minus "
		result += number_to_string(-1 * number)
	} else {
		result += smallnumbers[number]
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Too few arguments was given")
		return
	}

	input := os.Args[1]

	if input_number, err := strconv.Atoi(input); err == nil {
		fmt.Printf("%s is %s\n", input, number_to_string(input_number))
	} else {
		fmt.Println("Something went wrong with the int conversion -", err)
	}
}
