package walletdb

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/wallet"
	"github.com/stretchr/testify/assert"
)

func SaveAndVerify(t *testing.T) {
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
	repo := New(db)
	mnemonic := mnemonic.NewRandom()
	name := "test-wallet"
	wallet := wallet.New(mnemonic, "", name, 0, 0, "")
	if err := repo.Save(wallet); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	got, err := repo.Get(wallet.Name, "")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	assert.Equal(t, mnemonic.Words, got.Mnemonic.Words)
	assert.False(t, got.CreatedAt.IsZero(), "CreatedAt should be preserved")
	assert.True(t, got.CreatedAt.Equal(wallet.CreatedAt), "CreatedAt should match original")
}

func TestCreatedAtPersistence(t *testing.T) {
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
	repo := New(db)
	mnemonic := mnemonic.NewRandom()
	name := "test-wallet-timestamp"
	originalWallet := wallet.New(mnemonic, "", name, 5, 3, "")
	customTime := time.Date(2023, 1, 15, 10, 30, 45, 0, time.UTC)
	originalWallet.CreatedAt = customTime
	if err := repo.Save(originalWallet); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	retrievedWallet, err := repo.Get(name, "")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	assert.Equal(t, customTime, retrievedWallet.CreatedAt, "CreatedAt should be exactly preserved")
	assert.Equal(t, originalWallet.NextReceiveIndex, retrievedWallet.NextReceiveIndex)
	assert.Equal(t, originalWallet.NextChangeIndex, retrievedWallet.NextChangeIndex)
}

func CatchNotFound(t *testing.T) {
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
	repo := New(db)
	_, err = repo.Get("unknown-wallet", "")
	assert.EqualError(t, err, "wallet not found")
}
