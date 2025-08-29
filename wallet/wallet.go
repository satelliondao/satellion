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
	return child, nil
}

func (w *Wallet) NextIndex() uint32 {
	return w.NextChangeIndex
}

func (w *Wallet) ReceiveAddress() (*Address, error) {
	standart, err := w.RootKey.Derive(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return nil, err
	}
	coin, err := standart.Derive(hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}
	account, err := coin.Derive(hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}
	change, err := account.Derive(0)
	if err != nil {
		return nil, err
	}
	receive, err := change.Derive(w.NextReceiveIndex)
	if err != nil {
		return nil, err
	}
	addr, err := receive.Address(&chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	a := NewAddress(addr.EncodeAddress(), false, w.NextReceiveIndex)
	return a, nil
}

func (w *Wallet) NewReceiveAddress() (*Address, error) {
	w.NextReceiveIndex++
	return w.ReceiveAddress()
}
