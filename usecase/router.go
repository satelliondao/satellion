package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"time"

	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/persistence"
	"github.com/satelliondao/satellion/ports"
)

type Router struct {
	walletRepo *persistence.HDWalletRepo
}

func NewRouter() *Router {
	return &Router{
		walletRepo: persistence.NewHDWalletRepo(),
	}
}

func genNewWallet() *ports.HDWallet {
	// Generate a master private key (in real implementation, this would use BIP32/BIP39)
	masterPrivateKey := "0x" + hex.EncodeToString([]byte(fmt.Sprintf("master_key_%d", time.Now().UnixNano())))
	// Generate master address from private key
	masterAddress := DeriveAddressFromPrivateKey(masterPrivateKey)
	mnemonic := mnemonic.NewRandom()
	return &ports.HDWallet{
		MasterPrivateKey: masterPrivateKey,
		MasterPublicKey:  "0x" + hex.EncodeToString([]byte("master_public_key")),
		MasterAddress:    masterAddress,
		Mnemonic:         mnemonic,
		NextIndex:        0,
		UsedIndexes:      []uint32{},
	}
}

func createHDWalletFromSeed(mnemonic *mnemonic.Mnemonic) (*ports.HDWallet, error) {
	// Validate seed phrase (simplified)
	if len(mnemonic.Words) != 12 {
		return nil, fmt.Errorf("seed phrase must be 12 words")
	}
	// Generate master private key from seed phrase (simplified)
	hash := sha256.Sum256([]byte(mnemonic.String()))
	masterPrivateKey := "0x" + hex.EncodeToString(hash[:])
	masterAddress := DeriveAddressFromPrivateKey(masterPrivateKey)
	return &ports.HDWallet{
		MasterPrivateKey: masterPrivateKey,
		MasterPublicKey:  "0x" + hex.EncodeToString([]byte("master_public_key")),
		MasterAddress:    masterAddress,
		Mnemonic:         mnemonic,
		NextIndex:        0,
		UsedIndexes:      []uint32{},
	}, nil
}
