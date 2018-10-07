package main

import "fmt"
import "os"
import "strconv"

var (
	smallnumbers = map[int]string{
		0: "zero",
		1: "one",
		2: "two",
		3: "three",
		4: "four",
		5: "five",
		6: "six",
		7: "seven",
		8: "eight",
		9: "nine",
	}
)

func number_to_string(number int) string {
	return smallnumbers[number]
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
