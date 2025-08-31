package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/walletdb"
	"github.com/stretchr/testify/assert"
)

func setupTestWalletService(t *testing.T) (*WalletService, string, func()) {
	tmpDir, err := os.MkdirTemp("", "satellion-wallet-service-test-")
	if err != nil {
		t.Fatalf("mkdtemp failed: %v", err)
	}
	dbPath := filepath.Join(tmpDir, "test-wallets.db")
	db, err := walletdb.Connect(dbPath)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	repo := walletdb.New(db)
	service := NewWalletService(repo)
	cleanup := func() {
		db.Close()
		os.RemoveAll(tmpDir)
	}
	return service, tmpDir, cleanup
}

func TestWalletService_AddWallet_SetsCreatedAt(t *testing.T) {
	service, _, cleanup := setupTestWalletService(t)
	defer cleanup()
	words := []string{
		"abandon", "abandon", "abandon", "abandon", "abandon", "abandon",
		"abandon", "abandon", "abandon", "abandon", "abandon", "about",
	}
	testMnemonic := mnemonic.New(words)
	beforeAdd := time.Now()
	err := service.AddWallet("test-wallet", testMnemonic, "")
	afterAdd := time.Now()
	assert.NoError(t, err)
	retrievedWallet, err := service.walletRepo.Get("test-wallet", "")
	assert.NoError(t, err)
	assert.False(t, retrievedWallet.CreatedAt.IsZero(), "CreatedAt should be set")
	assert.True(t, retrievedWallet.CreatedAt.After(beforeAdd) || retrievedWallet.CreatedAt.Equal(beforeAdd))
	assert.True(t, retrievedWallet.CreatedAt.Before(afterAdd) || retrievedWallet.CreatedAt.Equal(afterAdd))
}

func TestWalletService_Unlock_ValidatesPassphrase(t *testing.T) {
	service, _, cleanup := setupTestWalletService(t)
	defer cleanup()
	words := []string{
		"abandon", "abandon", "abandon", "abandon", "abandon", "abandon",
		"abandon", "abandon", "abandon", "abandon", "abandon", "about",
	}
	testMnemonic := mnemonic.New(words)
	err := service.AddWallet("test-wallet", testMnemonic, "correct-passphrase")
	assert.NoError(t, err)
	err = service.Unlock("correct-passphrase")
	assert.NoError(t, err)
	err = service.Unlock("wrong-passphrase")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid passphrase")
}
