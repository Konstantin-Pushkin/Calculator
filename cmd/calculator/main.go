package main

import (
	"bufio"
	"calculator/pkg/calculator"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Print("Enter an expression: ")

	reader := bufio.NewReader(os.Stdin)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	result, err := calculator.Calc(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(52)
		return
	}

	fmt.Printf("Result: %v\n", result)
}
