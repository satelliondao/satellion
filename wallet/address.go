package wallet

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type Address struct {
	PubKey           *btcec.PublicKey
	Address          *btcutil.AddressTaproot
	Change           bool
	DeriviationIndex uint32
}

func NewAddress(
	pubKey *btcec.PublicKey,
	change bool,
	deriviationIndex uint32,
) *Address {
	// Compute the taproot output key (tweaked public key)
	outputKey, err := computeTaprootOutputKey(pubKey)
	if err != nil {
		panic(fmt.Errorf("failed to compute taproot output key: %w", err))
	}
	// This produces addresses starting with "bc1p" for mainnet
	taprootAddr, err := btcutil.NewAddressTaproot(outputKey, &chaincfg.MainNetParams)
	if err != nil {
		panic(fmt.Errorf("failed to create taproot address: %w", err))
	}
	return &Address{
		PubKey:           pubKey,
		Address:          taprootAddr,
		Change:           change,
		DeriviationIndex: deriviationIndex,
	}
}

// DeriveTaprootScriptPubKey generates a P2TR (Pay-to-Taproot) scriptPubKey
// following BIP 86 derivation path: m/86'/0'/0'/change/index
// Returns a 34-byte script: OP_1 <32-byte-taproot-output-key>
func (a *Address) DeriveTaprootScriptPubKey() ([]byte, error) {
	// Compute the taproot output key (tweaked public key)
	outputKey, err := computeTaprootOutputKey(a.PubKey)
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
func computeTaprootOutputKey(pubKey *btcec.PublicKey) ([]byte, error) {
	// Step 1: Serialize the internal public key to 32 bytes
	// SerializePubKey automatically ensures the y-coordinate is even (BIP 340 requirement)
	// This converts the public key to the schnorr format: 32-byte x-coordinate only
	internalKeyBytes := schnorr.SerializePubKey(pubKey)
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
