package mnemonic

import (
	"crypto/rand"
	"math/big"
	"strings"

	mnemonic "github.com/satelliondao/satellion/mnemonic/wordlist"
)

type Mnemonic struct {
	Words []string
}

func New(words []string) *Mnemonic {
	return &Mnemonic{
		Words: words,
	}
}

func NewRandom() *Mnemonic {
	out := make([]string, 12)
	max := big.NewInt(int64(len(mnemonic.EnWordList)))
	for i := 0; i < 12; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic("failed to generate random index")
		}
		out[i] = mnemonic.EnWordList[n.Int64()]
	}

	return &Mnemonic{
		Words: out,
	}
}

func (m *Mnemonic) String() string {
	return strings.Join(m.Words, " ")
}