package wallet

import (
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/satelliondao/satellion/mnemonic"
)

type Wallet struct {
	RootKey *hdkeychain.ExtendedKey
	nextIndex        uint32
	usedIndexes      []uint32
}

func New(mnemonic *mnemonic.Mnemonic) *Wallet {
	rootKey, _ := hdkeychain.NewMaster(mnemonic.Bytes(), &chaincfg.MainNetParams)
	return &Wallet{RootKey: rootKey}
}

