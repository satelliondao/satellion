package wallet

import (
	"fmt"
	"testing"

	"github.com/satelliondao/satellion/mnemonic"
)

func TestNewWallet(t *testing.T) {
	mnemonic := mnemonic.NewRandom()
	fmt.Printf("%+v\n", mnemonic)
	wallet := New(mnemonic)
	fmt.Printf("%+v\n", wallet)
	// ASSETT RROT KEY
	rootKey := wallet.RootKey
	fmt.Printf("%+v\n", rootKey)
}
