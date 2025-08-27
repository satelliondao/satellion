package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"time"

	"github.com/lightninglabs/neutrino/headerfs"
	"github.com/satelliondao/satellion/cfg"
	"github.com/satelliondao/satellion/chain"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/ports"
	"github.com/satelliondao/satellion/walletdb"
)

type Router struct {
	WalletRepo ports.WalletRepository
	Chain      *chain.Chain
	Config     *cfg.Config
}

func NewRouter() *Router {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	path := filepath.Join(home, ".satellion", "wallets.db")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		fmt.Println("failed to prepare wallets db dir:", err)
		os.Exit(1)
	}
	db, err := walletdb.Connect(path)
	if err != nil {
		fmt.Println("failed to open wallets db:", err)
		os.Exit(1)
	}
	loaded, _ := cfg.Load()
	repo := walletdb.NewWalletRepository(db)
	return &Router{WalletRepo: repo, Config: loaded}
}

func genNewWallet() *ports.HDWallet {
	// Generate a master private key (in real implementation, this would use BIP32/BIP39)
	masterPrivateKey := "0x" + hex.EncodeToString([]byte(fmt.Sprintf("master_key_%d", time.Now().UnixNano())))
	// Generate master address from private key
	mnemonic := mnemonic.NewRandom()
	return &ports.HDWallet{
		MasterPrivateKey: masterPrivateKey,
		MasterPublicKey:  "0x" + hex.EncodeToString([]byte("master_public_key")),
		MasterAddress:    "",
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
	return &ports.HDWallet{
		MasterPrivateKey: masterPrivateKey,
		MasterPublicKey:  "0x" + hex.EncodeToString([]byte("master_public_key")),
		MasterAddress:    "",
		Mnemonic:         mnemonic,
		NextIndex:        0,
		UsedIndexes:      []uint32{},
	}, nil
}

// UI Router integration helpers
func (r *Router) StartChain() error {
	if r.Chain != nil {
		return nil
	}
	if r.Config == nil {
		loaded, err := cfg.Load()
		if err != nil {
			return err
		}
		r.Config = loaded
	}
	r.Chain = chain.NewChain(r.Config)
	return r.Chain.Start()
}

func (r *Router) StopChain() error {
	if r.Chain == nil {
		return nil
	}
	err := r.Chain.Stop()
	r.Chain = nil
	return err
}

func (r *Router) BestBlock() (*headerfs.BlockStamp, int, error) {
	if r.Chain == nil {
		return nil, 0, fmt.Errorf("chain not started")
	}
	stamp, err := r.Chain.BestBlock()
	if err != nil {
		return nil, 0, err
	}
	return stamp, int(r.Chain.ConnectedCount()), nil
}

func (r *Router) MinPeers() int {
	if r.Config == nil || r.Config.MinPeers == 0 {
		return 5
	}
	return r.Config.MinPeers
}

// AddWallet saves the mnemonic under the provided wallet name.
func (r *Router) AddWallet(name string, m *mnemonic.Mnemonic) error {
	if name == "" || m == nil {
		return fmt.Errorf("invalid wallet data")
	}
	return r.WalletRepo.Add(name, m)
}
