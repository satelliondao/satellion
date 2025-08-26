package wallet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/satelliondao/satellion/mnemonic"
	wdb "github.com/satelliondao/satellion/walletdb"
	"github.com/stretchr/testify/assert"
)

func TestNewWallet(t *testing.T) {
	m := mnemonic.NewRandom()
	w := New(m)
	if w == nil {
		t.Fatalf("wallet should not be nil")
	}
}

func TestStoreWallet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "satellion-db-")
	if err != nil {
		t.Fatalf("mkdtemp failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := wdb.Connect(dbPath)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer db.Close()

	repo := wdb.NewWalletRepository(db)
	m := mnemonic.NewRandom()
	name := "test-wallet"

	if err := repo.Save(name, m); err != nil {
		t.Fatalf("store wallet failed: %v", err)
	}
	entity, err := repo.Get(name)
	if err != nil {
		t.Fatalf("get wallet failed: %v", err)
	}
	assert.Equal(t, entity.Words, m.Words)
}
