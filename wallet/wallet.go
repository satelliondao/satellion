package wallet

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
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
		Lock:             lock,
	}
	if lock == "" {
		seedH := sha256.Sum256(seed)
		w.Lock = hex.EncodeToString(seedH[:])
	}
	return w
}

func (w *Wallet) ReceiveAddress() (*Address, error) {
	return w.deriveTaprootAddress(0, w.NextReceiveIndex)
}

func (w *Wallet) ChangeAddress() (*Address, error) {
	return w.deriveTaprootAddress(1, w.NextChangeIndex)
}

func (w *Wallet) NewReceiveAddress() (*Address, error) {
	w.NextReceiveIndex++
	return w.ReceiveAddress()
}

func (w *Wallet) NewChangeAddress() (*Address, error) {
	w.NextChangeIndex++
	return w.ChangeAddress()
}

func (w *Wallet) deriveReceiveKeyPair(change, index uint32) (*btcec.PublicKey, *btcec.PrivateKey, error) {
	purpose, err := w.RootKey.Derive(hdkeychain.HardenedKeyStart + 86)
	if err != nil {
		return nil, nil, err
	}
	coin, err := purpose.Derive(hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, nil, err
	}
	account, err := coin.Derive(hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive account: %w", err)
	}
	// Derive change level (0 for receive, 1 for change)
	changePath, err := account.Derive(change)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive change path: %w", err)
	}
	// Derive the final address key at the specified index
	extendedKey, err := changePath.Derive(index)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive address key: %w", err)
	}
	// Step 6: Extract the public key from the derived extended key
	pubKey, err := extendedKey.ECPubKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get public key: %w", err)
	}
	privateKey, err := extendedKey.ECPrivKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get private key: %w", err)
	}
	return pubKey, privateKey, nil
}

// deriveTaprootAddress generates a BIP 86 taproot address
// following BIP 86 derivation path: m/86'/0'/0'/change/index
// Returns an Address struct with the bech32m-encoded taproot address (bc1p...)
func (w *Wallet) deriveTaprootAddress(change uint32, index uint32) (*Address, error) {
	pubKey, _, err := w.deriveReceiveKeyPair(change, index)
	if err != nil {
		return nil, fmt.Errorf("failed to derive receive key pair: %w", err)
	}
	return NewAddress(pubKey, change == 1, index), nil
}
