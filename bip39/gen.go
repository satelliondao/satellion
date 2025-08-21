package bip39

import (
	"crypto/rand"
	"math/big"
	"os"
	"strings"
)

func GenMnemonic() string {
	data, err := os.ReadFile("bip39/english.txt")
	if err != nil {
		panic("Bip39 word list not found")
	}
	raw := strings.Split(string(data), "\n")
	words := make([]string, 0, len(raw))
	for _, w := range raw {
		w = strings.TrimSpace(w)
		if w != "" {
			words = append(words, w)
		}
	}
	if len(words) == 0 {
		panic("Bip39 word list is empty")
	}
	out := make([]string, 12)
	max := big.NewInt(int64(len(words)))
	for i := 0; i < 12; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic("failed to generate random index")
		}
		out[i] = words[n.Int64()]
	}
	return strings.Join(out, " ")
}