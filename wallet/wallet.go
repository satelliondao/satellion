package wallet

import (
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/satelliondao/satellion/mnemonic"
)

type Wallet struct {
	RootKey   *hdkeychain.ExtendedKey
	nextIndex uint32
}

func New(mnemonic *mnemonic.Mnemonic) *Wallet {
	rootKey, _ := hdkeychain.NewMaster(mnemonic.Bytes(), &chaincfg.MainNetParams)
	return &Wallet{RootKey: rootKey, nextIndex: 0}
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
