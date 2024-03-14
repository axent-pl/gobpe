package main

import (
	"fmt"

	"github.com/axent-pl/gobpe/preprocessor"
	"github.com/axent-pl/gobpe/tokenizer"
)

const DEFAULT_SPLIT_PATTERN = `(?i:[sdmt]|ll|ve|re)|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]|\s+`
const DEFAULT_PREPROCESSOR_PATTERN = `([;.]|^)(\p{L})`
const DEFAULT_PREPROCESSOR_REPLACEMENT = `$1 $2`

func main() {
	links, err := LoadLinksFromFile(`data/url/bajki.txt`)
	if err != nil {
		panic(err)
	}
	paths, err := Download(links, `data/txt/`)
	if err != nil {
		panic(err)
	}

	tokenizr := tokenizer.New(
		tokenizer.WithSplitPattern(DEFAULT_SPLIT_PATTERN),
		tokenizer.WithPreprocessor(preprocessor.New(DEFAULT_PREPROCESSOR_PATTERN, DEFAULT_PREPROCESSOR_REPLACEMENT)),
	)

	for _, filePath := range paths {
		text := MustReadFile(filePath)
		tokenizr.LoadText(text)
	}

	tokenizr.Fit(1000, 1000)

	tokenizerBytes, _ := tokenizr.Serialize()
	MustWriteToFile("params.json", tokenizerBytes)

	tokens := tokenizr.StringTokens()
	for tokenId, tokenString := range tokens {
		fmt.Printf("%v: %v\n", tokenId, tokenString)
	}

	text2 := MustReadFile("sample.txt")
	encoded := tokenizr.Encode(text2)
	decoded := string(tokenizr.Decode(encoded))
	fmt.Println(decoded)
	fmt.Printf("Length before encoding: %v, Length after encoding: %v\n", len(text2), len(encoded))

	tokenizerBytes2 := MustReadFile("params.json")
	tokenizer2 := tokenizer.Tokenizer{}
	tokenizer2.Deserialize(tokenizerBytes2)
	encoded2 := tokenizr.Encode(text2)
	decoded2 := string(tokenizr.Decode(encoded2))
	fmt.Println(decoded2)
	fmt.Printf("Length before encoding: %v, Length after encoding: %v\n", len(text2), len(encoded2))
}
