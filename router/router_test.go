package router

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcutil/gcs"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/neutrino/headerfs"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/neutrino"
	"github.com/satelliondao/satellion/ports"
	"github.com/satelliondao/satellion/wallet"
	"github.com/satelliondao/satellion/walletdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChain struct {
	mock.Mock
}

// Compile-time check to ensure MockChain implements ports.ChainService
var _ ports.ChainService = (*MockChain)(nil)

func (m *MockChain) BestBlock() (*headerfs.BlockStamp, error) {
	args := m.Called()
	return args.Get(0).(*headerfs.BlockStamp), args.Error(1)
}

func (m *MockChain) GetBlockHash(height int64) (*chainhash.Hash, error) {
	args := m.Called(height)
	return args.Get(0).(*chainhash.Hash), args.Error(1)
}

func (m *MockChain) GetBlockHeader(hash *chainhash.Hash) (*wire.BlockHeader, error) {
	args := m.Called(hash)
	return args.Get(0).(*wire.BlockHeader), args.Error(1)
}

func (m *MockChain) GetCFilter(hash chainhash.Hash) (*gcs.Filter, error) {
	args := m.Called(hash)
	return args.Get(0).(*gcs.Filter), args.Error(1)
}

func (m *MockChain) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockChain) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockChain) ConnectedCount() int32 {
	args := m.Called()
	return args.Get(0).(int32)
}

func (m *MockChain) IsSynced() bool {
	args := m.Called()
	return args.Bool(0)
}

func setupTestRouter(t *testing.T) (*Router, string, func()) {
	tmpDir, err := os.MkdirTemp("", "satellion-router-test-")
	if err != nil {
		t.Fatalf("mkdtemp failed: %v", err)
	}
	dbPath := filepath.Join(tmpDir, "test-wallets.db")
	db, err := walletdb.Connect(dbPath)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	repo := walletdb.New(db)
	cfg := &config.Config{MinPeers: 1}
	router := &Router{
		WalletRepo: repo,
		Config:     cfg,
	}
	cleanup := func() {
		db.Close()
		os.RemoveAll(tmpDir)
	}
	return router, tmpDir, cleanup
}

func TestGetWalletBalance_ChainNotStarted(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()
	balance, err := router.GetWalletBalance("test")
	assert.Error(t, err)
	assert.Equal(t, uint64(0), balance)
	assert.Contains(t, err.Error(), "chain not started")
}

func TestGetWalletBalance_NoActiveWallet(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()
	router.Chain = &neutrino.ChainService{}
	balance, err := router.GetWalletBalance("test")
	assert.Error(t, err)
	assert.Equal(t, uint64(0), balance)
	assert.Contains(t, err.Error(), "wallet not found")
}

func TestGetWalletBalance_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()
	words := []string{
		"abandon", "abandon", "abandon", "abandon", "abandon", "abandon",
		"abandon", "abandon", "abandon", "abandon", "abandon", "about",
	}
	testMnemonic := mnemonic.New(words)
	testWallet := wallet.New(&testMnemonic, "", "test-wallet", 0, 0, "")
	testWallet.CreatedAt = time.Now().Add(-24 * time.Hour)
	err := router.WalletRepo.Save(testWallet)
	assert.NoError(t, err)
	err = router.WalletRepo.SetDefault("test-wallet")
	assert.NoError(t, err)
	// Skip this test since it requires actual chain integration
	// The balance scanning functionality is tested in the wallet package
	t.Skip("Skipping integration test - balance scanning tested in wallet package")
}

func TestAddWallet_SetsCreatedAt(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()
	words := []string{
		"abandon", "abandon", "abandon", "abandon", "abandon", "abandon",
		"abandon", "abandon", "abandon", "abandon", "abandon", "about",
	}
	testMnemonic := mnemonic.New(words)
	beforeAdd := time.Now()
	err := router.AddWallet("test-wallet", testMnemonic, "")
	afterAdd := time.Now()
	assert.NoError(t, err)
	retrievedWallet, err := router.WalletRepo.Get("test-wallet", "")
	assert.NoError(t, err)
	assert.False(t, retrievedWallet.CreatedAt.IsZero(), "CreatedAt should be set")
	assert.True(t, retrievedWallet.CreatedAt.After(beforeAdd) || retrievedWallet.CreatedAt.Equal(beforeAdd))
	assert.True(t, retrievedWallet.CreatedAt.Before(afterAdd) || retrievedWallet.CreatedAt.Equal(afterAdd))
}

func TestUnlock_ValidatesPassphrase(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()
	words := []string{
		"abandon", "abandon", "abandon", "abandon", "abandon", "abandon",
		"abandon", "abandon", "abandon", "abandon", "abandon", "about",
	}
	testMnemonic := mnemonic.New(words)
	err := router.AddWallet("test-wallet", testMnemonic, "correct-passphrase")
	assert.NoError(t, err)
	err = router.Unlock("correct-passphrase")
	assert.NoError(t, err)
	err = router.Unlock("wrong-passphrase")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid passphrase")
}
