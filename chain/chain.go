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
	"github.com/fatih/color"
	"github.com/lightninglabs/neutrino"
	"github.com/lightninglabs/neutrino/headerfs"
	"github.com/satelliondao/satellion/cfg"
	"github.com/satelliondao/satellion/utils/term"
	"github.com/satelliondao/satellion/walletdb"
)

type Chain struct {
	chainService *neutrino.ChainService
	cfg          *cfg.Config
}

func NewChain(
	cfg *cfg.Config,
) *Chain {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dataDir := filepath.Join(home, ".satellion", "neutrino", "mainnet")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		log.Fatal("failed to create data dir: ", err)
	}
	db, err := walletdb.Connect(dataDir)
	if err != nil {
		log.Fatal("failed to open neutrino db: ", err)
	}
	chainService, err := neutrino.NewChainService(neutrino.Config{
		DataDir:     dataDir,
		Database:    db,
		ChainParams: chaincfg.MainNetParams,
		AddPeers:    cfg.Peers,
	})
	if err != nil {
		log.Fatal(`failed to create chain service: `, err)
	}
	return &Chain{
		chainService: chainService,
		cfg:          cfg,
	}
}

func (c *Chain) BestBlock() (*headerfs.BlockStamp, error) {
	return c.chainService.BestBlock()
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
				term.PrintfInline("best block error: %v\n", err)
				continue
			}
			term.PrintfInline("best height=%d time=%s peers=%d", stamp.Height, stamp.Timestamp.UTC().Format(time.RFC3339), c.chainService.ConnectedCount())

			if int(c.chainService.ConnectedCount()) >= c.cfg.MinPeers {
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
			term.Newline()
			fmt.Println("shutting down...")
			return nil
		}
	}
}

func (c *Chain) Stop() error {
	return c.chainService.Stop()
}
