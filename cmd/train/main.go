package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/axent-pl/gobpe/preprocessor"
	"github.com/axent-pl/gobpe/tokenizer"
)

const DEFAULT_SPLIT_PATTERN = `(?i:[sdmt]|ll|ve|re)|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]|\s+`
const DEFAULT_PREPROCESSOR_PATTERN = `([;.]|^)(\p{L})`
const DEFAULT_PREPROCESSOR_REPLACEMENT = `$1 $2`

func main() {
	var paramsFilePath string
	var splitPattern string
	var preprocessorPattern string
	var preprocessorReplacement string
	var maxIterations int
	var maxToken int

	flag.StringVar(&paramsFilePath, "params", "params.json", "Path to params file")
	flag.StringVar(&splitPattern, "splitPattern", DEFAULT_SPLIT_PATTERN, "Token split pattern")
	flag.StringVar(&preprocessorPattern, "preprocessorPattern", DEFAULT_PREPROCESSOR_PATTERN, "Preprocessor pattern")
	flag.StringVar(&preprocessorReplacement, "preprocessorReplacement", DEFAULT_PREPROCESSOR_REPLACEMENT, "Preprocessor replacement")
	flag.IntVar(&maxIterations, "maxIterations", 1000, "Max number of iterationsations")
	flag.IntVar(&maxToken, "maxToken", 1000, "Max number of tokens")
	flag.Parse()

	tokenizr := tokenizer.New(
		tokenizer.WithSplitPattern(splitPattern),
		tokenizer.WithPreprocessor(preprocessor.New(preprocessorPattern, preprocessorReplacement)),
	)

	text, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Error reading stdin: %s", err)
		os.Exit(1)
	}

	tokenizr.LoadText(text)
	tokenizr.Fit(maxIterations, maxToken)
	tokenizerBytes, _ := tokenizr.Serialize()

	err = os.WriteFile(paramsFilePath, tokenizerBytes, 0600)
	if err != nil {
		panic(err)
	}
}
