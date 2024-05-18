package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/axent-pl/gobpe/tokenizer"
)

func main() {
	var paramsFilePath string

	flag.StringVar(&paramsFilePath, "params", "params.json", "Path to params file")

	tokenizer := tokenizer.New(tokenizer.FromSerialized(paramsFilePath))

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Error reading stdin: %s", err)
		os.Exit(1)
	}

	encodedString := tokenizer.DecodeFromString(data)
	fmt.Println(string(encodedString))
}