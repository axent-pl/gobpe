package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

const MaxSimpleToken = 255

type Pair struct {
	Value [2]int
}

func (p Pair) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]int{p.Value[0], p.Value[1]})
}

func (p *Pair) UnmarshalJSON(data []byte) error {
	var arr [2]int
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	p.Value = arr
	return nil
}

type Preprocessor struct {
	Pattern     regexp.Regexp `json:"pattern"`
	Replacement []byte        `json:"replacement"`
}

func (p Preprocessor) Process(text []byte) []byte {
	return p.Pattern.ReplaceAll(text, p.Replacement)
}

type Tokenizer struct {
	Tokens            [][]int       `json:"-"`
	TextPreprocessor  Preprocessor  `json:"preprocessor"`
	SplitPattern      regexp.Regexp `json:"split_regexp"`
	LastToken         int           `json:"last_token"`
	ReplacementKeys   []Pair        `json:"replacement_keys"`
	ReplacementValues []int         `json:"replacement_values"`
}

func NewTokenizer(text []byte, split *regexp.Regexp) *Tokenizer {
	t := &Tokenizer{
		SplitPattern: *split,
		LastToken:    MaxSimpleToken,
	}
	byteSequences := split.FindAll(text, -1)
	t.tokensFromBytes(byteSequences)
	return t
}

func (t *Tokenizer) Serialize() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Tokenizer) Deserialize(data []byte) error {
	var tc Tokenizer
	if err := json.Unmarshal(data, &tc); err != nil {
		return err
	}
	t.LastToken = tc.LastToken
	t.ReplacementKeys = tc.ReplacementKeys
	t.ReplacementValues = tc.ReplacementValues
	return nil
}

func (t *Tokenizer) tokensFromBytes(byteSequences [][]byte) {
	t.Tokens = make([][]int, 0, len(byteSequences))
	for _, bytes := range byteSequences {
		intSlice := make([]int, 0, len(bytes))
		for _, b := range bytes {
			intSlice = append(intSlice, int(b))
		}
		t.Tokens = append(t.Tokens, intSlice)
	}
}

func (t *Tokenizer) Fit(maxIterations int, maxTokenValue int) {
	for i := 0; i < maxIterations; i++ {
		pair, count := t.FindMostFrequentPair()
		if count < 2 {
			break
		}
		t.LastToken++
		if t.LastToken > maxTokenValue {
			break
		}
		t.ReplacePair(pair, t.LastToken)
	}
}

func (t *Tokenizer) ReplacePair(pair Pair, value int) {
	t.ReplacementKeys = append(t.ReplacementKeys, pair)
	t.ReplacementValues = append(t.ReplacementValues, value)
	for i, tokenSeq := range t.Tokens {
		newSeq := make([]int, 0, len(tokenSeq))
		for j := 0; j < len(tokenSeq); {
			if j+1 < len(tokenSeq) && tokenSeq[j] == pair.Value[0] && tokenSeq[j+1] == pair.Value[1] {
				newSeq = append(newSeq, value)
				j += 2
			} else {
				newSeq = append(newSeq, tokenSeq[j])
				j++
			}
		}
		t.Tokens[i] = newSeq
	}
}

func (t *Tokenizer) FindMostFrequentPair() (Pair, int) {
	var counter = map[Pair]int{}
	var maxCount int = 0
	var maxKey Pair = Pair{}

	for i := 0; i < len(t.Tokens); i++ {
		for c := 0; c < len(t.Tokens[i])-1; c++ {
			key := Pair{Value: [2]int{t.Tokens[i][c], t.Tokens[i][c+1]}}
			count, exists := counter[key]
			if exists {
				counter[key] = count + 1
				if count+1 > maxCount {
					maxKey = key
					maxCount = count + 1
				}
			} else {
				counter[key] = 1
			}
		}
	}
	return maxKey, maxCount
}

func (t *Tokenizer) Encode(text []byte) []int {
	var encoded = make([]int, 0)
	for _, b := range text {
		encoded = append(encoded, int(b))
	}
	for i := 0; i < len(t.ReplacementKeys); i++ {
		encodedReplaced := make([]int, 0, len(encoded))
		for c := 0; c < len(encoded); {
			if c+1 < len(encoded) && encoded[c] == t.ReplacementKeys[i].Value[0] && encoded[c+1] == t.ReplacementKeys[i].Value[1] {
				encodedReplaced = append(encodedReplaced, t.ReplacementValues[i])
				c += 2
			} else {
				encodedReplaced = append(encodedReplaced, encoded[c])
				c += 1
			}
		}
		encoded = encodedReplaced
	}
	return encoded
}

func (t *Tokenizer) Decode(tokens []int) []byte {
	var decoded = make([]int, 0, len(tokens))
	decoded = tokens

	for i := len(t.ReplacementKeys) - 1; i >= 0; i-- {
		decodedReplaced := make([]int, 0, len(decoded))
		for c := 0; c < len(decoded); c++ {
			if decoded[c] > MaxSimpleToken {
				idx := decoded[c] - MaxSimpleToken - 1
				decodedReplaced = append(decodedReplaced, t.ReplacementKeys[idx].Value[0])
				decodedReplaced = append(decodedReplaced, t.ReplacementKeys[idx].Value[1])
			} else {
				decodedReplaced = append(decodedReplaced, decoded[c])
			}
		}
		decoded = decodedReplaced
	}

	var decodedBytes = make([]byte, 0, len(decoded))
	for c := 0; c < len(decoded); c++ {
		decodedBytes = append(decodedBytes, byte(decoded[c]))
	}
	return decodedBytes
}

func (t *Tokenizer) StringTokens() map[int]string {
	var tokens map[int]string = make(map[int]string)
	var result map[int][]byte = make(map[int][]byte)
	for i := 0; i < len(t.ReplacementKeys); i++ {
		tokenId := t.ReplacementValues[i]
		result[tokenId] = make([]byte, 0)
		for j := 0; j < 2; j++ {
			if t.ReplacementKeys[i].Value[j] > MaxSimpleToken {
				idx := t.ReplacementKeys[i].Value[j]
				result[tokenId] = append(result[tokenId], result[idx]...)
			} else {
				result[tokenId] = append(result[tokenId], byte(t.ReplacementKeys[i].Value[j]))
			}
		}
		tokens[tokenId] = string(result[tokenId])
	}
	return tokens
}

func main() {
	preprocessor := Preprocessor{
		Pattern:     *regexp.MustCompile(`([[:punct:]]|^)(\p{L})`),
		Replacement: []byte(`$1 $2`),
	}
	split := regexp.MustCompile(`(^\p{L}+| \p{L}+| [0-9]+|[[:punct:]]|[[:space:]]+)`)

	text := MustReadFile("lorem.txt")
	textProcessed := preprocessor.Process(text)

	tokenizer := NewTokenizer(textProcessed, split)
	tokenizer.Fit(1000, 1000)
	tokenizerBytes, _ := tokenizer.Serialize()
	MustWriteToFile("params.json", tokenizerBytes)

	tokens := tokenizer.StringTokens()
	for tokenId, tokenString := range tokens {
		fmt.Printf("%v: %v\n", tokenId, tokenString)
	}

	text2 := MustReadFile("latin.txt")
	encoded := tokenizer.Encode(text2)
	decoded := string(tokenizer.Decode(encoded))
	fmt.Println(decoded)
	fmt.Printf("Length before encoding: %v, Length after encoding: %v\n", len(text2), len(encoded))

	tokenizerBytes2 := MustReadFile("params.json")
	tokenizer2 := Tokenizer{}
	tokenizer2.Deserialize(tokenizerBytes2)
	encoded2 := tokenizer.Encode(text2)
	decoded2 := string(tokenizer.Decode(encoded2))
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
