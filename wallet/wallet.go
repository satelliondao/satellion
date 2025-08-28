package wallet

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/satelliondao/satellion/mnemonic"
)

type Wallet struct {
	Mnemonic  *mnemonic.Mnemonic
	RootKey   *hdkeychain.ExtendedKey
	nextIndex uint32
	Name      string
	IsDefault bool
}

func New(mnemonic *mnemonic.Mnemonic, name string, nextIndex uint32) *Wallet {
	seed := mnemonic.Seed("")
	rootKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		panic(fmt.Sprintf("failed to create root key: %v", err))
	}
	return &Wallet{
		RootKey:   rootKey,
		nextIndex: nextIndex,
		Name:      name,
		Mnemonic:  mnemonic,
	}
}

func (w *Wallet) DeriveChild(index uint32) (*hdkeychain.ExtendedKey, error) {
	child, err := w.RootKey.Derive(index)
	if err != nil {
		return nil, err
	}
	w.nextIndex++
	return child, nil
}

func (w *Wallet) NextIndex() uint32 {
	return w.nextIndex
}
