package wallet

import (
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcutil/gcs"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/neutrino/headerfs"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const passphrase = "test"

func TestNewBalanceScanner(t *testing.T) {
	mockChain := &MockChainService{}
	scanner := NewBalanceService(mockChain)
	assert.NotNil(t, scanner)
	assert.Equal(t, mockChain, scanner.chain)
}

func TestScanWalletBalance_WalletCreatedAtZero(t *testing.T) {
	mockChain := &MockChainService{}
	scanner := NewBalanceService(mockChain)
	words := []string{
		"abandon", "abandon", "abandon", "abandon", "abandon", "abandon",
		"abandon", "abandon", "abandon", "abandon", "abandon", "about",
	}
	mnemonic := mnemonic.New(words)
	wallet := New(&mnemonic, passphrase, "test")
	wallet.CreatedAt = time.Time{} // Zero time
	balance, err := scanner.ScanWalletBalanceInfo(wallet)
	assert.Error(t, err)
	assert.Equal(t, uint64(0), balance)
	assert.Contains(t, err.Error(), "wallet creation time not set")
}

func TestScanWalletBalance_ChainError(t *testing.T) {
	mockChain := &MockChainService{}
	scanner := NewBalanceService(mockChain)
	words := []string{
		"abandon", "abandon", "abandon", "abandon", "abandon", "abandon",
		"abandon", "abandon", "abandon", "abandon", "abandon", "about",
	}
	mnemonic := mnemonic.New(words)
	wallet := New(&mnemonic, passphrase, "test")
	wallet.CreatedAt = time.Now().Add(-24 * time.Hour)
	mockChain.On("BestBlock").Return((*headerfs.BlockStamp)(nil), assert.AnError)
	balance, err := scanner.ScanWalletBalanceInfo(wallet)
	assert.Error(t, err)
	assert.Equal(t, uint64(0), balance)
	assert.Contains(t, err.Error(), "failed to get best block")
}

func TestScanWalletBalance_Success(t *testing.T) {
	mockChain := &MockChainService{}
	scanner := NewBalanceService(mockChain)
	words := []string{
		"abandon", "abandon", "abandon", "abandon", "abandon", "abandon",
		"abandon", "abandon", "abandon", "abandon", "abandon", "about",
	}
	mnemonic := mnemonic.New(words)
	wallet := New(&mnemonic, passphrase, "test")
	wallet.CreatedAt = time.Now().Add(-24 * time.Hour)
	bestBlock := &headerfs.BlockStamp{
		Height:    100,
		Hash:      chainhash.Hash{},
		Timestamp: time.Now(),
	}
	mockChain.On("BestBlock").Return(bestBlock, nil)
	blockHash := &chainhash.Hash{}
	blockHeader := &wire.BlockHeader{
		Timestamp: wallet.CreatedAt.Add(time.Hour),
	}
	mockChain.On("GetBlockHash", mock.AnythingOfType("int64")).Return(blockHash, nil)
	mockChain.On("GetBlockHeader", blockHash).Return(blockHeader, nil)
	var key [16]byte
	emptyFilter, _ := gcs.BuildGCSFilter(0, 0, key, [][]byte{})
	mockChain.On("GetCFilter", mock.AnythingOfType("chainhash.Hash"), wire.GCSFilterRegular).Return(emptyFilter, nil)
	balance, err := scanner.ScanWalletBalanceInfo(wallet)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), balance)
}

func TestFindBlockHeightFromTime_BinarySearch(t *testing.T) {
	mockChain := &MockChainService{}
	scanner := NewBalanceService(mockChain)
	targetTime := time.Now()
	blockHash := &chainhash.Hash{}
	mockChain.On("GetBlockHash", mock.AnythingOfType("int64")).Return(blockHash, nil)
	mockChain.On("GetBlockHeader", blockHash).Return(&wire.BlockHeader{
		Timestamp: targetTime.Add(-time.Hour),
	}, nil).Once()
	mockChain.On("GetBlockHeader", blockHash).Return(&wire.BlockHeader{
		Timestamp: targetTime.Add(time.Hour),
	}, nil)
	height, err := scanner.findBlockHeightFromTime(targetTime, 100)
	assert.NoError(t, err)
	assert.True(t, height >= 0)
}

func TestGenerateAllAddresses(t *testing.T) {
	words := []string{
		"abandon", "abandon", "abandon", "abandon", "abandon", "abandon",
		"abandon", "abandon", "abandon", "abandon", "abandon", "about",
	}
	mnemonic := mnemonic.New(words)
	wallet := New(&mnemonic, passphrase, "test")
	mockChain := &MockChainService{}
	scanner := NewBalanceService(mockChain)
	addresses, err := scanner.generateAllAddresses(wallet)
	assert.NoError(t, err)
	expectedCount := (3 + 1) * 2 // (max receive index + 1) * 2 (receive + change)
	assert.Equal(t, expectedCount, len(addresses))
	receiveCount := 0
	changeCount := 0
	for _, addr := range addresses {
		if addr.Change {
			changeCount++
		} else {
			receiveCount++
		}
	}
	assert.Equal(t, 4, receiveCount)
	assert.Equal(t, 4, changeCount)
}

func TestGenerateAllAddresses_ZeroIndices(t *testing.T) {
	words := []string{
		"abandon", "abandon", "abandon", "abandon", "abandon", "abandon",
		"abandon", "abandon", "abandon", "abandon", "abandon", "about",
	}
	mnemonic := mnemonic.New(words)
	wallet := New(&mnemonic, passphrase, "test")
	mockChain := &MockChainService{}
	scanner := NewBalanceService(mockChain)
	addresses, err := scanner.generateAllAddresses(wallet)
	assert.NoError(t, err)
	expectedCount := 21 * 2 // Default 20 addresses + index 0, both receive and change
	assert.Equal(t, expectedCount, len(addresses))
}

func TestAddressesToScripts(t *testing.T) {
	words := []string{
		"abandon", "abandon", "abandon", "abandon", "abandon", "abandon",
		"abandon", "abandon", "abandon", "abandon", "abandon", "about",
	}
	mnemonic := mnemonic.New(words)
	wallet := New(&mnemonic, passphrase, "test")
	mockChain := &MockChainService{}
	scanner := NewBalanceService(mockChain)
	addresses, err := scanner.generateAllAddresses(wallet)
	assert.NoError(t, err)
	scripts, err := scanner.addressesToScripts(addresses)
	assert.NoError(t, err)
	assert.Equal(t, len(addresses), len(scripts))
	for _, script := range scripts {
		assert.True(t, len(script) > 0)
	}
}

func TestScanBlockForAddresses_NoMatches(t *testing.T) {
	mockChain := &MockChainService{}
	scanner := NewBalanceService(mockChain)
	blockHash := &chainhash.Hash{}
	mockChain.On("GetBlockHash", int64(100)).Return(blockHash, nil)
	var key [16]byte
	emptyFilter, _ := gcs.BuildGCSFilter(0, 0, key, [][]byte{})
	mockChain.On("GetCFilter", *blockHash, wire.GCSFilterRegular).Return(emptyFilter, nil)
	scripts := [][]byte{
		{0x51, 0x20, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}
	matches, err := scanner.scanBlockForAddresses(100, scripts)
	assert.NoError(t, err)
	assert.Equal(t, 0, matches)
}

func TestScanBlockForAddresses_GetBlockHashError(t *testing.T) {
	mockChain := &MockChainService{}
	scanner := NewBalanceService(mockChain)
	mockChain.On("GetBlockHash", int64(100)).Return((*chainhash.Hash)(nil), assert.AnError)
	scripts := [][]byte{{0x51, 0x20}}
	matches, err := scanner.scanBlockForAddresses(100, scripts)
	assert.Error(t, err)
	assert.Equal(t, 0, matches)
	assert.Contains(t, err.Error(), "failed to get block hash")
}

func TestScanBlockForAddresses_GetCFilterError(t *testing.T) {
	mockChain := &MockChainService{}
	scanner := NewBalanceService(mockChain)
	blockHash := &chainhash.Hash{}
	mockChain.On("GetBlockHash", int64(100)).Return(blockHash, nil)
	mockChain.On("GetCFilter", *blockHash, wire.GCSFilterRegular).Return((*gcs.Filter)(nil), assert.AnError)
	scripts := [][]byte{{0x51, 0x20}}
	matches, err := scanner.scanBlockForAddresses(100, scripts)
	assert.Error(t, err)
	assert.Equal(t, 0, matches)
	assert.Contains(t, err.Error(), "failed to get compact filter")
}
