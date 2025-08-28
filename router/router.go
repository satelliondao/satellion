package router

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lightninglabs/neutrino/headerfs"
	"github.com/satelliondao/satellion/chain"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/wallet"
	"github.com/satelliondao/satellion/walletdb"
)

type Router struct {
	WalletRepo *walletdb.WalletDB
	Chain      *chain.Chain
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
func (r *Router) AddWallet(name string, m mnemonic.Mnemonic) error {
	if name == "" {
		return fmt.Errorf("invalid wallet data")
	}
	err := r.WalletRepo.Save(wallet.New(&m, name, 0))
	if err != nil {
		return err
	}
	return r.WalletRepo.SetDefault(name)
}
