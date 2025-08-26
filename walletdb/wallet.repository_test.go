package walletdb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/satelliondao/satellion/mnemonic"
	"github.com/stretchr/testify/assert"
)

func TestWalletRepository_SaveAndGet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "satellion-db-")
	if err != nil {
		t.Fatalf("mkdtemp failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := Connect(dbPath)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer db.Close()
	repo := NewWalletRepository(db)
	m := mnemonic.NewRandom()
	name := "test-wallet"
	if err := repo.Save(name, m); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	got, err := repo.Get(name)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	assert.Equal(t, m.Words, got.Words)
}

func TestWalletRepository_Get_NotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "satellion-db-")
	if err != nil {
		t.Fatalf("mkdtemp failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := Connect(dbPath)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer db.Close()
	repo := NewWalletRepository(db)
	_, err = repo.Get("unknown-wallet")
	assert.EqualError(t, err, "wallet not found")
}
