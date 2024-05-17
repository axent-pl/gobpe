package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var sourcePath string
	flag.StringVar(&sourcePath, "src", "", "Path to train files")
	flag.Parse()

	paths, err := listFiles(sourcePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	for _, filePath := range paths {
		cleanFile(filePath)
	}
}
