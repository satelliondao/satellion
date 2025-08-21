package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"strings"
	"time"

	"github.com/satelliondao/satellion/bip39"
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
	seedPhrase := bip39.GenMnemonic()
	return &ports.HDWallet{
		MasterPrivateKey: masterPrivateKey,
		MasterPublicKey:  "0x" + hex.EncodeToString([]byte("master_public_key")),
		MasterAddress:    masterAddress,
		SeedPhrase:       seedPhrase,
		NextIndex:        0,
		UsedIndexes:      []uint32{},
	}
}

func createHDWalletFromSeed(seedPhrase string) (*ports.HDWallet, error) {
	// Validate seed phrase (simplified)
	if len(strings.Split(seedPhrase, " ")) != 12 {
		return nil, fmt.Errorf("seed phrase must be 12 words")
	}
	// Generate master private key from seed phrase (simplified)
	hash := sha256.Sum256([]byte(seedPhrase))
	masterPrivateKey := "0x" + hex.EncodeToString(hash[:])
	masterAddress := DeriveAddressFromPrivateKey(masterPrivateKey)
	return &ports.HDWallet{
		MasterPrivateKey: masterPrivateKey,
		MasterPublicKey:  "0x" + hex.EncodeToString([]byte("master_public_key")),
		MasterAddress:    masterAddress,
		SeedPhrase:       seedPhrase,
		NextIndex:        0,
		UsedIndexes:      []uint32{},
	}, nil
}

