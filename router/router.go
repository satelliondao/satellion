package router

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lightninglabs/neutrino/headerfs"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/neutrino"
	"github.com/satelliondao/satellion/wallet"
	"github.com/satelliondao/satellion/walletdb"
)

type Router struct {
	WalletRepo *walletdb.WalletDB
	Chain      *neutrino.ChainService
	Config     *config.Config
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
	loaded, _ := config.Load()
	repo := walletdb.New(db)
	return &Router{WalletRepo: repo, Config: loaded}
}

// UI Router integration helpers
func (r *Router) StartChain() error {
	if r.Chain != nil {
		return nil
	}
	if r.Config == nil {
		loaded, err := config.Load()
		if err != nil {
			return err
		}
		r.Config = loaded
	}
	r.Chain = neutrino.NewChainService(r.Config)
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

func (r *Router) SyncTimeoutMinutes() int {
	if r.Config == nil || r.Config.SyncTimeoutMinutes == 0 {
		return 30
	}
	return r.Config.SyncTimeoutMinutes
}

// AddWallet saves the mnemonic under the provided wallet name.
func (r *Router) AddWallet(name string, m mnemonic.Mnemonic, passphrase string) error {
	if name == "" {
		return fmt.Errorf("invalid wallet data")
	}
	model := wallet.New(&m, passphrase, "")
	model.Name = name
	model.CreatedAt = time.Now()

	err := r.WalletRepo.Save(model)
	if err != nil {
		return err
	}
	return r.WalletRepo.SetDefault(name)
}

// Unlock derives and stores the xprv for the active wallet using the provided 13th-word passphrase.
func (r *Router) Unlock(passphrase string) error {
	w, err := r.WalletRepo.GetActiveWallet(passphrase)
	if err != nil {
		return err
	}
	if w == nil || w.Mnemonic == nil {
		return fmt.Errorf("no active wallet")
	}
	hashSeed := sha256.Sum256(w.Mnemonic.Seed(passphrase))
	unlockKey := hex.EncodeToString(hashSeed[:])

	if w.Lock != unlockKey {
		return fmt.Errorf("invalid passphrase")
	}
	return nil
}

// GetWalletBalanceInfo scans for wallet balance and UTXO count using compact filters from creation time
func (r *Router) GetWalletBalanceInfo(passphrase string) (*neutrino.BalanceInfo, error) {
	if r.Chain == nil {
		return nil, fmt.Errorf("chain not started")
	}
	w, err := r.WalletRepo.GetActiveWallet(passphrase)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, fmt.Errorf("no active wallet")
	}
	balanceService := neutrino.NewBalanceService(r.Chain)
	return balanceService.ScanLedger(w)
}
