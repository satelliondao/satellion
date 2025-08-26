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
	walletRepo ports.WalletRepository
	chain      *chain.Chain
	cfgLoaded  *cfg.Config
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
	return &Router{walletRepo: walletdb.NewWalletRepository(db), cfgLoaded: loaded}
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
	if r.chain != nil {
		return nil
	}
	if r.cfgLoaded == nil {
		loaded, err := cfg.Load()
		if err != nil {
			return err
		}
		r.cfgLoaded = loaded
	}
	r.chain = chain.NewChain(r.cfgLoaded)
	return r.chain.Start()
}

func (r *Router) StopChain() error {
	if r.chain == nil {
		return nil
	}
	err := r.chain.Stop()
	r.chain = nil
	return err
}

func (r *Router) BestBlock() (*headerfs.BlockStamp, int, error) {
	if r.chain == nil {
		return nil, 0, fmt.Errorf("chain not started")
	}
	stamp, err := r.chain.BestBlock()
	if err != nil {
		return nil, 0, err
	}
	return stamp, int(r.chain.ConnectedCount()), nil
}

func (r *Router) MinPeers() int {
	if r.cfgLoaded == nil || r.cfgLoaded.MinPeers == 0 {
		return 5
	}
	return r.cfgLoaded.MinPeers
}

// AddWallet saves the mnemonic under the provided wallet name.
func (r *Router) AddWallet(name string, m *mnemonic.Mnemonic) error {
	if name == "" || m == nil {
		return fmt.Errorf("invalid wallet data")
	}
	return r.walletRepo.Add(name, m)
}
