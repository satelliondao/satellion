package chain

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	bdb "github.com/btcsuite/btcwallet/walletdb"
	"github.com/fatih/color"
	"github.com/lightninglabs/neutrino"
	"github.com/lightninglabs/neutrino/headerfs"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/walletdb"
)

type Chain struct {
	chainService *neutrino.ChainService
	config       *config.Config
	db           bdb.DB
}

func NewChain(
	config *config.Config,
) *Chain {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dataDir := filepath.Join(home, ".satellion", "neutrino", "mainnet")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		log.Fatal("failed to create data dir: ", err)
	}
	dbPath := filepath.Join(dataDir, "neutrino.db")
	db, err := walletdb.Connect(dbPath)
	if err != nil {
		log.Fatal("failed to open neutrino db: ", err)
	}
	chainService, err := neutrino.NewChainService(neutrino.Config{
		DataDir:     dataDir,
		Database:    db,
		ChainParams: chaincfg.MainNetParams,
		AddPeers:    config.Peers,
	})
	if err != nil {
		log.Fatal(`failed to create chain service: `, err)
	}
	return &Chain{
		chainService: chainService,
		config:       config,
		db:           db,
	}
}

func (c *Chain) BestBlock() (*headerfs.BlockStamp, error) {
	return c.chainService.BestBlock()
}

func (c *Chain) Start() error {
	return c.chainService.Start()
}

func (c *Chain) ConnectedCount() int32 {
	return c.chainService.ConnectedCount()
}

func (c *Chain) Sync() error {
	if err := c.chainService.Start(); err != nil {
		return err
	}
	defer func() { _ = c.chainService.Stop() }()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			stamp, err := c.chainService.BestBlock()
			if err != nil {
				continue
			}

			if int(c.chainService.ConnectedCount()) >= c.config.MinPeers {
				isCurrent := false
				type isCurrentCap interface{ IsCurrent() bool }
				if v, ok := interface{}(c.chainService).(isCurrentCap); ok {
					isCurrent = v.IsCurrent()
				} else {
					if time.Since(stamp.Timestamp) < 10*time.Minute {
						isCurrent = true
					}
				}
				if isCurrent {
					green := color.New(color.FgGreen)
					green.Println("\nSynchronization complete, head at height", stamp.Height)
					return nil
				}
			}
		case <-sigCh:
			fmt.Println("shutting down...")
			return nil
		}
	}
}

func (c *Chain) Stop() error {
	_ = c.chainService.Stop()
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
