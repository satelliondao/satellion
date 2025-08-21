package mnemonic

import (
	"fmt"
	"strings"

	mnemonic "github.com/satelliondao/satellion/mnemonic/wordlist"
)

type Validator struct {
	wordmap map[string]struct{}
}

func NewValidator() *Validator {
	wordmap  := make(map[string]struct{})
	for _, word := range mnemonic.EnWordList {
		wordmap[word] = struct{}{}
	}
	return &Validator{
		wordmap: wordmap,
	}
}

func (v *Validator) IsWord(word string) bool {
	word = strings.ToLower(word)
	_, ok := v.wordmap[word]
	return ok
}

func (v *Validator) Validate(mnemonic string) error {
	words := v.Normalize(mnemonic)
	if len(words) != 12 {
		return fmt.Errorf("mnemonic must be 12 words")
	}

	for _, word := range words {
		if !v.IsWord(word) {
			return fmt.Errorf("invalid word: %s", word)
		}
	}

	return nil
}

func (v *Validator) Normalize(mnemonic string) []string {
	words := strings.Split(strings.TrimSpace(mnemonic), " ")
	normalized := make([]string, 0, len(words))
	for _, word := range words {
		normalized = append(normalized, strings.ToLower(word))
	}
	return normalized
}