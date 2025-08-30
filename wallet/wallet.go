package wallet

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
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

func (w *Wallet) ReceiveScriptPubKey() ([]byte, error) {
	return w.DeriveTaprootScriptPubKey(0, w.NextReceiveIndex)
}

func (w *Wallet) ChangeScriptPubKey() ([]byte, error) {
	return w.DeriveTaprootScriptPubKey(1, w.NextChangeIndex)
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
	// Compute the taproot output key (tweaked public key)
	outputKey, err := computeTaprootOutputKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to compute taproot output key: %w", err)
	}
	// Create the taproot address using bech32m encoding
	// This produces addresses starting with "bc1p" for mainnet
	taprootAddr, err := btcutil.NewAddressTaproot(outputKey, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create taproot address: %w", err)
	}
	return NewAddress(taprootAddr.String(), change == 1, index), nil
}

// DeriveTaprootScriptPubKey generates a P2TR (Pay-to-Taproot) scriptPubKey
// following BIP 86 derivation path: m/86'/0'/0'/change/index
// Returns a 34-byte script: OP_1 <32-byte-taproot-output-key>
func (w *Wallet) DeriveTaprootScriptPubKey(change uint32, index uint32) ([]byte, error) {
	pubKey, _, err := w.deriveReceiveKeyPair(change, index)
	if err != nil {
		return nil, fmt.Errorf("failed to derive receive key pair: %w", err)
	}
	// Compute the taproot output key (tweaked public key)
	outputKey, err := computeTaprootOutputKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to compute taproot output key: %w", err)
	}
	// Construct the P2TR scriptPubKey
	// Format: OP_1 (0x51) + PUSH_32 (0x20) + 32-byte-output-key
	scriptPubKey := make([]byte, 34)
	scriptPubKey[0] = 0x51            // OP_1: version 1 witness program
	scriptPubKey[1] = 0x20            // Push 32 bytes (0x20 = 32 in hex)
	copy(scriptPubKey[2:], outputKey) // Copy the 32-byte taproot output key
	return scriptPubKey, nil
}

// computeTaprootOutputKey implements BIP 341 taproot output key computation.
// It takes an internal public key and computes the taproot output key using the formula:
// Q = P + int(hash_TapTweak(bytes(P))) * G
// where P is the internal key, Q is the output key, and G is the generator point.
func computeTaprootOutputKey(pubk *btcec.PublicKey) ([]byte, error) {
	// Step 1: Serialize the internal public key to 32 bytes
	// SerializePubKey automatically ensures the y-coordinate is even (BIP 340 requirement)
	// This converts the public key to the schnorr format: 32-byte x-coordinate only
	internalKeyBytes := schnorr.SerializePubKey(pubk)
	// Step 2: Compute the TapTweak hash
	// BIP 341 defines: hash_TapTweak(P) = tagged_hash("TapTweak", P)
	// This creates a deterministic tweak value from the internal key
	tapTweakHash := chainhash.TaggedHash(chainhash.TagTapTweak, internalKeyBytes)
	// Step 3: Parse the serialized key back to get the canonical version
	// This ensures we have the public key with an even y-coordinate
	// (required for schnorr signatures and taproot)
	evenInternalKey, err := schnorr.ParsePubKey(internalKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse internal key: %w", err)
	}
	// Step 4: Convert the 32-byte tweak hash to a scalar for elliptic curve operations
	// ModNScalar represents integers modulo the secp256k1 curve order
	tweakInt := new(btcec.ModNScalar)
	tweakInt.SetByteSlice(tapTweakHash[:])
	// Step 5: Compute the tweak point (t * G)
	// This multiplies the tweak scalar by the generator point G
	// JacobianPoint is used for efficient elliptic curve arithmetic
	tweakPoint := new(btcec.JacobianPoint)
	btcec.ScalarBaseMultNonConst(tweakInt, tweakPoint)
	tweakPoint.ToAffine() // Convert from Jacobian to affine coordinates
	// Step 6: Convert the internal key to Jacobian coordinates for point addition
	internalPoint := new(btcec.JacobianPoint)
	evenInternalKey.AsJacobian(internalPoint)
	// Step 7: Compute the final output point (P + t*G)
	// This is the core taproot computation: internal key + tweak point
	// AddNonConst performs elliptic curve point addition
	btcec.AddNonConst(internalPoint, tweakPoint, internalPoint)
	internalPoint.ToAffine() // Convert back to affine coordinates
	// Step 8: Create the final output public key and serialize it
	// This gives us the 32-byte taproot output key that goes into the scriptPubKey
	outputKey := btcec.NewPublicKey(&internalPoint.X, &internalPoint.Y)
	return schnorr.SerializePubKey(outputKey), nil
}
