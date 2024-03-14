package main

import (
	"fmt"
	"os"

	"github.com/axent-pl/gobpe/preprocessor"
	"github.com/axent-pl/gobpe/tokenizer"
)

func main() {
	tokenizr := tokenizer.New(
		tokenizer.WithSplitPattern(`(?i:[sdmt]|ll|ve|re)|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]|\s+`),
		tokenizer.WithPreprocessor(preprocessor.New(`([;.]|^)(\p{L})`, `$1 $2`)),
	)

	text := MustReadFile("lorem.txt")
	tokenizr.LoadText(text)
	tokenizr.Fit(1000, 1000)

	tokenizerBytes, _ := tokenizr.Serialize()
	MustWriteToFile("params.json", tokenizerBytes)

	tokens := tokenizr.StringTokens()
	for tokenId, tokenString := range tokens {
		fmt.Printf("%v: %v\n", tokenId, tokenString)
	}

	text2 := MustReadFile("latin.txt")
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

func MustWriteToFile(filename string, data []byte) {
	err := os.WriteFile(filename, data, 0600)
	if err != nil {
		panic(err)
	}
}

func MustReadFile(filepath string) []byte {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return bytes
}
