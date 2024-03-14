package preprocessor

import "regexp"

type Preprocessor struct {
	Pattern     regexp.Regexp `json:"pattern"`
	Replacement []byte        `json:"replacement"`
}

func New(pattern string, replacement string) *Preprocessor {
	preprocessor := &Preprocessor{
		Pattern:     *regexp.MustCompile(pattern),
		Replacement: []byte(replacement),
	}
	return preprocessor
}

func (p Preprocessor) Process(text []byte) []byte {
	return p.Pattern.ReplaceAll(text, p.Replacement)
}
