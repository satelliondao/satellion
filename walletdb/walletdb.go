package walletdb

import (
	"os"
	"time"

	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb"
)

var defaultDBTimeout = 10 * time.Second

func NewWalletDB() (walletdb.DB, error) {
	return walletdb.Create("bdb", "satellion.db", true, 60*time.Second, false)
}

func Open(path string) (walletdb.DB, error) {
	_, statErr := os.Stat(path)
	if os.IsNotExist(statErr) {
		return walletdb.Create("bdb", path, true, 60*time.Second, false)
	}
	return walletdb.Open("bdb", path, true, 60*time.Second, false)
}