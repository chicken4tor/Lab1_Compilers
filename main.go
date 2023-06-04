package main

import (
	"fmt"
	"os"
)

func main() {
	input, err := os.Open("ioFiles/inputFile.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	output, err := os.Create("ioFiles/outputFile.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	LexIt(input, output)
}
