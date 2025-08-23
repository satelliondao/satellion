package mnemonic

import (
	"crypto/rand"
	"math/big"
	"strings"

	mnemonic "github.com/satelliondao/satellion/mnemonic/wordlist"
)

const defaultMnemonicWordCount = 12

type Mnemonic struct {
	Words []string
}

func New(words []string) *Mnemonic {
	return &Mnemonic{
		Words: words,
	}
}

func NewRandom() *Mnemonic {
	out := make([]string, defaultMnemonicWordCount)
	max := big.NewInt(int64(len(mnemonic.EnWordList)))
	for i := 0; i < defaultMnemonicWordCount; i++ {
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

func (m *Mnemonic) Bytes() []byte {
	return []byte(m.String())
}