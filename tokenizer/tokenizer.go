package tokenizer

import (
	"encoding/json"
	"regexp"

	"github.com/axent-pl/gobpe/preprocessor"
)

const MaxSimpleToken = 255

type Tokenizer struct {
	Tokens            [][]int                    `json:"-"`
	Preprocessor      *preprocessor.Preprocessor `json:"preprocessor"`
	SplitPattern      regexp.Regexp              `json:"split_regexp"`
	LastToken         int                        `json:"last_token"`
	ReplacementKeys   []Pair                     `json:"replacement_keys"`
	ReplacementValues []int                      `json:"replacement_values"`
}

func WithPreprocessor(p *preprocessor.Preprocessor) func(*Tokenizer) {
	return func(t *Tokenizer) {
		t.Preprocessor = p
	}
}

func WithSplitPattern(pattern string) func(*Tokenizer) {
	return func(t *Tokenizer) {
		t.SplitPattern = *regexp.MustCompile(pattern)
	}
}

func New(options ...func(*Tokenizer)) *Tokenizer {
	t := &Tokenizer{
		LastToken: MaxSimpleToken,
		Tokens:    make([][]int, 0),
	}
	for _, option := range options {
		option(t)
	}
	return t
}

func (t *Tokenizer) LoadText(text []byte) error {
	if t.Preprocessor != nil {
		text = t.Preprocessor.Process(text)
	}
	sequences := t.SplitPattern.FindAll(text, -1)
	for _, seq := range sequences {
		intSequence := make([]int, 0, len(seq))
		for _, b := range seq {
			intSequence = append(intSequence, int(b))
		}
		t.Tokens = append(t.Tokens, intSequence)
	}
	return nil
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
	if t.Preprocessor != nil {
		text = t.Preprocessor.Process(text)
	}
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
