package wallet

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/satelliondao/satellion/mnemonic"
)

type Wallet struct {
	Mnemonic         *mnemonic.Mnemonic
	RootKey          *hdkeychain.ExtendedKey
	NextChangeIndex  uint32
	NextReceiveIndex uint32
	Name             string
	IsDefault        bool
	Passphrase       string
	Lock             string
}

func New(
	mnemonic *mnemonic.Mnemonic,
	passphrase string,
	name string,
	nextChangeIndex uint32,
	nextReceiveIndex uint32,
	lock string,
) *Wallet {
	seed := mnemonic.Seed(passphrase)
	rootKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		panic(fmt.Sprintf("failed to create root key: %v", err))
	}
	w := &Wallet{
		RootKey:          rootKey,
		NextChangeIndex:  nextChangeIndex,
		NextReceiveIndex: nextReceiveIndex,
		Name:             name,
		Mnemonic:         mnemonic,
		Passphrase:       passphrase,
		Lock:             lock,
	}
	if lock == "" {
		seedH := sha256.Sum256(seed)
		w.Lock = hex.EncodeToString(seedH[:])
	}
	return w
}

func (w *Wallet) DeriveChild(index uint32) (*hdkeychain.ExtendedKey, error) {
	child, err := w.RootKey.Derive(index)
	if err != nil {
		return nil, err
	}
	w.NextChangeIndex++
	return child, nil
}

func (w *Wallet) NextIndex() uint32 {
	return w.NextChangeIndex
}
