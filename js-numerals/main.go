package main

import "fmt"
import "os"
import "strconv"

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
)

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}

func number_to_string(number int) string {
	if number > 19 {
		_, remainder := divmod(number, 10)
		if remainder != 0 {
			return fmt.Sprintf("%s-%s", tens[number-remainder], smallnumbers[remainder])
		} else {
			return fmt.Sprintf("%s", tens[number-remainder])
		}
	} else {
		return smallnumbers[number]
	}
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
