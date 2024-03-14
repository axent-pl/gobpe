package tokenizer

import "encoding/json"

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
