package walletdb

import (
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb"
)

var (
	defaultDBTimeout = 10 * time.Second
	walletStoreKey   = []byte("wallets")
)

type DB struct {
	db walletdb.DB
}

func Connect(dataDir string) (walletdb.DB, error) {
	if dataDir != "" {
		if err := os.MkdirAll(filepath.Dir(dataDir), 0o755); err != nil {
			return nil, err
		}
		return openOrCreate(dataDir)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dataDir = filepath.Join(home, ".satellion", "neutrino", "mainnet")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(dataDir, "neutrino.db")
	return openOrCreate(path)
}

func openOrCreate(path string) (walletdb.DB, error) {
	_, statErr := os.Stat(path)
	if os.IsNotExist(statErr) {
		return walletdb.Create("bdb", path, true, defaultDBTimeout, false)
	}
	return walletdb.Open("bdb", path, true, defaultDBTimeout, false)
}
