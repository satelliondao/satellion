package walletdb

import (
	"os"
	"path/filepath"
	"testing"

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
	wallet := wallet.New(mnemonic, "", name, 0, 0)
	if err := repo.Save(wallet); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	got, err := repo.Get(wallet.Name, "")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	assert.Equal(t, mnemonic.Words, got.Mnemonic.Words)
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
