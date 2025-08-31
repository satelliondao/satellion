package neutrino

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/btcsuite/btcd/btcutil/gcs"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	bdb "github.com/btcsuite/btcwallet/walletdb"
	"github.com/fatih/color"
	"github.com/lightninglabs/neutrino"
	"github.com/lightninglabs/neutrino/headerfs"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/ports"
	"github.com/satelliondao/satellion/walletdb"
)

type ChainService struct {
	neutrino *neutrino.ChainService
	config   *config.Config
	db       bdb.DB
}

var _ ports.ChainService = (*ChainService)(nil)

func NewChainService(
	config *config.Config,
) *ChainService {
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
	neutrino, err := neutrino.NewChainService(neutrino.Config{
		DataDir:     dataDir,
		Database:    db,
		ChainParams: chaincfg.MainNetParams,
		AddPeers:    config.Peers,
	})
	if err != nil {
		log.Fatal(`failed to create chain service: `, err)
	}
	return &ChainService{
		neutrino: neutrino,
		config:   config,
		db:       db,
	}
}

func (c *ChainService) GetCFilter(hash chainhash.Hash) (*gcs.Filter, error) {
	return c.neutrino.GetCFilter(hash, wire.GCSFilterRegular)
}

func (c *ChainService) GetBlockHash(height int64) (*chainhash.Hash, error) {
	return c.neutrino.GetBlockHash(height)
}

func (c *ChainService) BestBlock() (*headerfs.BlockStamp, error) {
	return c.neutrino.BestBlock()
}

func (c *ChainService) Start() error {
	return c.neutrino.Start()
}

func (c *ChainService) ConnectedCount() int32 {
	return c.neutrino.ConnectedCount()
}

func (c *ChainService) Sync() error {
	if err := c.neutrino.Start(); err != nil {
		return err
	}
	defer func() { _ = c.neutrino.Stop() }()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			stamp, err := c.neutrino.BestBlock()
			if err != nil {
				continue
			}

			if int(c.neutrino.ConnectedCount()) >= c.config.MinPeers {
				isCurrent := false
				type isCurrentCap interface{ IsCurrent() bool }
				if v, ok := interface{}(c.neutrino).(isCurrentCap); ok {
					isCurrent = v.IsCurrent()
				} else {
					if time.Since(stamp.Timestamp) < time.Duration(c.config.SyncTimeoutMinutes)*time.Minute {
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

func (c *ChainService) IsSynced() bool {
	return c.neutrino.IsCurrent()
}

func (c *ChainService) Stop() error {
	_ = c.neutrino.Stop()
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *ChainService) GetBlockHeader(hash *chainhash.Hash) (*wire.BlockHeader, error) {
	return c.neutrino.GetBlockHeader(hash)
}
