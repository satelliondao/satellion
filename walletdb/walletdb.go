package walletdb

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb"
)

var (
	defaultDBTimeout = 10 * time.Second
	walletStoreKey       = []byte("wallets")
)

type DB struct {
	db walletdb.DB
}

func Connect() (walletdb.DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dataDir := filepath.Join(home, ".satellion", "neutrino", "mainnet")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		log.Fatal("failed to create data dir: ", err)
	}
	path := filepath.Join(dataDir, "neutrino.db")
	db, err := openOrCreate(path)
	if err != nil {
		log.Fatal("failed to open neutrino db: ", err)
	}

	return db, nil
}


func openOrCreate(path string) (walletdb.DB, error) {
	_, statErr := os.Stat(path)
	if os.IsNotExist(statErr) {
		return walletdb.Create("bdb", path, true, 60*time.Second, false)
	}
	return walletdb.Open("bdb", path, true, 60*time.Second, false)
}
