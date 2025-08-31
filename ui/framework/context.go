package framework

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/neutrino"
	"github.com/satelliondao/satellion/service"
	"github.com/satelliondao/satellion/walletdb"
)

type AppContext struct {
	Passphrase    string
	WalletService *service.WalletService
	ChainService  *neutrino.Chain
	Config        *config.Config
	WalletRepo    *walletdb.WalletDB
}

func NewContext() (*AppContext, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	path := filepath.Join(home, ".satellion", "wallets.db")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("failed to prepare wallets db dir: %w", err)
	}
	db, err := walletdb.Connect(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open wallets db: %w", err)
	}
	loaded, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	repo := walletdb.New(db)
	walletService := service.NewWalletService(repo)
	chainService, err := neutrino.NewChain(loaded)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain service: %w", err)
	}
	return &AppContext{
		WalletService: walletService,
		ChainService:  chainService,
		Config:        loaded,
		WalletRepo:    repo,
	}, nil
}

func (ctx *AppContext) Cleanup() {
	if ctx.ChainService != nil {
		ctx.ChainService.Stop()
	}
}
