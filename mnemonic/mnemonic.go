package mnemonic

import (
	"crypto/rand"
	"crypto/sha512"
	"math/big"
	"strings"

	"golang.org/x/crypto/pbkdf2"

	"github.com/satelliondao/satellion/mnemonic/wordlist"
)

const wordCount = 12
const wordlistCount = 2048

type Mnemonic struct {
	Words []string
}

func New(words []string) Mnemonic {
	return Mnemonic{
		Words: words,
	}
}

func NewRandom() *Mnemonic {
	out := make([]string, wordCount)
	if wordlistCount != len(wordlist.EnWordList) {
		panic("wordlist count mismatch")
	}
	max := big.NewInt(wordlistCount)
	for i := range wordCount {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic("failed to generate random index")
		}
		out[i] = wordlist.EnWordList[n.Int64()]
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

func (m *Mnemonic) Seed(passphrase string) []byte {
	return pbkdf2.Key([]byte(m.String()), []byte("mnemonic"+passphrase), 2048, 64, sha512.New)
}
