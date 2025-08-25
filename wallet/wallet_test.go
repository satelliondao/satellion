package wallet

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	bdb "github.com/btcsuite/btcwallet/walletdb"
	"github.com/satelliondao/satellion/mnemonic"
	wdb "github.com/satelliondao/satellion/walletdb"
)

func TestNewWallet(t *testing.T) {
	m := mnemonic.NewRandom()
	w := New(m)
	if w == nil { t.Fatalf("wallet should not be nil") }
}

func TestStoreWallet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "satellion-db-")
	if err != nil { t.Fatalf("mkdtemp failed: %v", err) }
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := wdb.Connect(dbPath)
	if err != nil { t.Fatalf("connect failed: %v", err) }
	defer db.Close()

	repo := wdb.NewWalletRepository(db)
	m := mnemonic.NewRandom()
	name := "test-wallet"

	if err := repo.Save(name, m); err != nil { t.Fatalf("store wallet failed: %v", err) }
	err = db.View(func(tx bdb.ReadTx) error {
		b := tx.ReadBucket([]byte("wallet_" + name))
		if b == nil { t.Fatalf("dedicated bucket not found") }
		raw := b.Get([]byte("wallet"))
		if len(raw) == 0 { t.Fatalf("wallet key not found") }
		var entity wdb.WalletEntity
		if uErr := json.Unmarshal(raw, &entity); uErr != nil { t.Fatalf("unmarshal failed: %v", uErr) }
		if entity.Name != name || entity.Mnemonic != m.String() { t.Fatalf("persisted wallet mismatch") }
		return nil
	}, func() {})
	
	if err != nil { t.Fatalf("db view failed: %v", err) }
}
