package main

import "fmt"
import "os"
import "strconv"

func number_to_string(number int) string {
	return ""
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
