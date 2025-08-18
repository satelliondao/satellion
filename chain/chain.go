package chain

import (
	"log"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg"
	n "github.com/lightninglabs/neutrino"
	"github.com/lightninglabs/neutrino/headerfs"
	"github.com/satelliondao/satellion/db"
)

type ChainService struct {
	Chain *n.ChainService
}

func NewChainService(
) *ChainService {
	panic("use NewChainServiceWithPeers instead")
}

func NewChainServiceWithPeers(
	peers []string,
) *ChainService {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dataDir := filepath.Join(home, ".satellion", "neutrino", "mainnet")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		log.Fatal("failed to create data dir: ", err)
	}

	walletDBPath := filepath.Join(dataDir, "neutrino.db")
	boltDB, err := db.OpenOrCreate(walletDBPath)
	if err != nil {
		log.Fatal("failed to open neutrino db: ", err)
	}

	chainService, err := n.NewChainService(n.Config{
		DataDir:     dataDir,
		Database:    boltDB,
		ChainParams: chaincfg.MainNetParams,
		AddPeers:    peers,
	})
	if err != nil {
		log.Fatal(`failed to create chain service: `, err)
	}
	return &ChainService{
		Chain: chainService,
	}
}

func (c *ChainService) BestBlock() (*headerfs.BlockStamp, error) {
	return c.Chain.BestBlock()
}

func (c *ChainService) Start() error {
	return c.Chain.Start()
}

func (c *ChainService) Stop() error {
	return c.Chain.Stop()
}

