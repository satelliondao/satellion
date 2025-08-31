package neutrino

import (
	"fmt"
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
	"github.com/lightninglabs/neutrino"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/ports"
	"github.com/satelliondao/satellion/walletdb"
)

type Chain struct {
	neutrino *neutrino.ChainService
	config   *config.Config
	db       bdb.DB
}

var _ ports.Chain = (*Chain)(nil)

func NewChain(config *config.Config) (*Chain, error) {
	var s = &Chain{config: config}
	if s.config == nil {
		loaded, err := config.Load()
		if err != nil {
			return nil, err
		}
		s.config = loaded
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dataDir := filepath.Join(home, ".satellion", "neutrino", "mainnet")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create data dir: %w", err)
	}
	dbPath := filepath.Join(dataDir, "neutrino.db")

	s.db, err = walletdb.Connect(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open neutrino db: %w", err)
	}

	s.neutrino, err = neutrino.NewChainService(neutrino.Config{
		DataDir:     dataDir,
		Database:    s.db,
		ChainParams: chaincfg.MainNetParams,
		AddPeers:    s.config.Peers,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create chain service: %w", err)
	}
	return s, nil
}

func (s *Chain) BestBlock() (*ports.BlockInfo, error) {
	stamp, err := s.neutrino.BestBlock()
	if err != nil {
		return nil, err
	}
	peers := int(s.neutrino.ConnectedCount())
	return &ports.BlockInfo{
		BlockStamp: stamp,
		Peers:      peers,
	}, nil
}

func (s *Chain) GetCFilter(hash chainhash.Hash) (*gcs.Filter, error) {
	return s.neutrino.GetCFilter(hash, wire.GCSFilterRegular)
}

func (s *Chain) GetBlockHash(height int64) (*chainhash.Hash, error) {
	return s.neutrino.GetBlockHash(height)
}

func (s *Chain) GetBlockHeader(hash *chainhash.Hash) (*wire.BlockHeader, error) {
	return s.neutrino.GetBlockHeader(hash)
}

func (s *Chain) ConnectedCount() int32 {
	return s.neutrino.ConnectedCount()
}

func (s *Chain) Syncronize() error {
	if err := s.neutrino.Start(); err != nil {
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			stamp, err := s.neutrino.BestBlock()
			if err != nil {
				continue
			}
			if int(s.neutrino.ConnectedCount()) >= s.config.MinPeers {
				isCurrent := false
				type isCurrentCap interface{ IsCurrent() bool }
				if v, ok := interface{}(s.neutrino).(isCurrentCap); ok {
					isCurrent = v.IsCurrent()
				} else {
					if time.Since(stamp.Timestamp) < time.Duration(s.config.SyncTimeoutMinutes)*time.Minute {
						isCurrent = true
					}
				}
				if isCurrent {
					return nil
				}
			}
		case <-sigCh:
			fmt.Println("shutting down...")
			return nil
		}
	}
}

func (s *Chain) IsSynced() bool {
	return s.neutrino.IsCurrent()
}

func (s *Chain) Stop() {
	if s.neutrino != nil {
		s.neutrino.Stop()
	}
}
