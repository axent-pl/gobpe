package main

import (
	"reflect"
	"testing"
)

func TestTokenizer_ReplacePair(t *testing.T) {
	type fields struct {
		Tokens            [][]int
		MaxToken          int
		ReplacementKeys   []Pair
		ReplacementValues []int
	}
	type args struct {
		pair  Pair
		value int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantTokens [][]int
	}{
		{
			name: "Replace single occurrence",
			fields: fields{
				Tokens: [][]int{
					{1, 2, 3},
					{4, 5, 6},
				},
				MaxToken:          255,
				ReplacementKeys:   nil,
				ReplacementValues: nil,
			},
			args: args{
				pair:  Pair{Value: [2]int{2, 3}},
				value: 256,
			},
			wantTokens: [][]int{
				{1, 256},
				{4, 5, 6},
			},
		},
		{
			name: "Replace multiple occurrences",
			fields: fields{
				Tokens: [][]int{
					{1, 2, 3, 2, 3},
					{2, 3, 4},
				},
				MaxToken:          255,
				ReplacementKeys:   nil,
				ReplacementValues: nil,
			},
			args: args{
				pair:  Pair{Value: [2]int{2, 3}},
				value: 256,
			},
			wantTokens: [][]int{
				{1, 256, 256},
				{256, 4},
			},
		},
		{
			name: "No replacement due to no occurrence",
			fields: fields{
				Tokens: [][]int{
					{1, 2, 4},
					{5, 6, 7},
				},
				MaxToken:          255,
				ReplacementKeys:   nil,
				ReplacementValues: nil,
			},
			args: args{
				pair:  Pair{Value: [2]int{3, 4}},
				value: 256,
			},
			wantTokens: [][]int{
				{1, 2, 4},
				{5, 6, 7},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tokenizer{
				Tokens:            tt.fields.Tokens,
				LastToken:         tt.fields.MaxToken,
				ReplacementKeys:   tt.fields.ReplacementKeys,
				ReplacementValues: tt.fields.ReplacementValues,
			}
			tr.ReplacePair(tt.args.pair, tt.args.value)
			for i, tokenSeq := range tr.Tokens {
				if !reflect.DeepEqual(tokenSeq, tt.wantTokens[i]) {
					t.Errorf("ReplacePair() Tokens[%d] = %v, want %v", i, tokenSeq, tt.wantTokens[i])
				}
			}
		})
	}
}

func TestTokenizer_FindMostFrequentPair(t *testing.T) {
	type fields struct {
		Tokens            [][]int
		MaxToken          int
		ReplacementKeys   []Pair
		ReplacementValues []int
	}
	tests := []struct {
		name      string
		fields    fields
		wantPair  Pair
		wantCount int
	}{
		{
			name: "Simple find",
			fields: fields{
				Tokens: [][]int{
					{1, 2, 3},
					{2, 3, 4},
				},
				MaxToken:          255,
				ReplacementKeys:   nil,
				ReplacementValues: nil,
			},
			wantPair:  Pair{Value: [2]int{2, 3}},
			wantCount: 2,
		},
		{
			name: "No find",
			fields: fields{
				Tokens: [][]int{
					{1, 2, 3},
					{5, 6, 7},
				},
				MaxToken:          255,
				ReplacementKeys:   nil,
				ReplacementValues: nil,
			},
			wantPair:  Pair{},
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tokenizer{
				Tokens:            tt.fields.Tokens,
				LastToken:         tt.fields.MaxToken,
				ReplacementKeys:   tt.fields.ReplacementKeys,
				ReplacementValues: tt.fields.ReplacementValues,
			}
			got, got1 := tr.FindMostFrequentPair()
			if !reflect.DeepEqual(got, tt.wantPair) {
				t.Errorf("Tokenizer.FindMostFrequentPair() got = %v, want %v", got, tt.wantPair)
			}
			if got1 != tt.wantCount {
				t.Errorf("Tokenizer.FindMostFrequentPair() got1 = %v, want %v", got1, tt.wantCount)
			}
		})
	}
}
