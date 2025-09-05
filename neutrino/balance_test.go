package neutrino

import (
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcutil/gcs"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/neutrino/headerfs"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const passphrase = "test"

var words = []string{
	"abandon", "abandon", "abandon", "abandon", "abandon", "abandon",
	"abandon", "abandon", "abandon", "abandon", "abandon", "about",
}
var seed = mnemonic.New(words)

func setupTest() (*MockChainService, *BalanceService) {
	mockChain := &MockChainService{}
	scanner := NewBalance(mockChain)
	return mockChain, scanner
}

func TestNewBalanceService(t *testing.T) {
	chain, scanner := setupTest()
	assert.NotNil(t, scanner)
	assert.Equal(t, chain, scanner.chain)
}

func TestScanWalletBalance_WalletCreatedAtZero(t *testing.T) {
	_, scanner := setupTest()
	w := wallet.New(&seed, passphrase, "test")
	w.CreatedAt = time.Time{} // Zero time
	balance, err := scanner.ScanLedger(w)
	assert.Error(t, err)
	assert.Nil(t, balance)
	assert.Contains(t, err.Error(), "wallet creation time not set")
}

func TestScanWalletBalance_ChainError(t *testing.T) {
	chain, scanner := setupTest()
	w := wallet.New(&seed, passphrase, "test")
	w.CreatedAt = time.Now().Add(-24 * time.Hour)
	chain.On("BestBlock").Return((*headerfs.BlockStamp)(nil), assert.AnError)
	balance, err := scanner.ScanLedger(w)
	assert.Error(t, err)
	assert.Nil(t, balance)
	assert.Contains(t, err.Error(), "failed to get best block")
}

func TestScanWalletBalance_Success(t *testing.T) {
	chain, scanner := setupTest()
	w := wallet.New(&seed, passphrase, "test")
	w.CreatedAt = time.Now().Add(-24 * time.Hour)
	bestBlock := &headerfs.BlockStamp{
		Height:    100,
		Hash:      chainhash.Hash{},
		Timestamp: time.Now(),
	}
	chain.On("BestBlock").Return(bestBlock, nil)
	blockHash := &chainhash.Hash{}
	blockHeader := &wire.BlockHeader{
		Timestamp: w.CreatedAt.Add(time.Hour),
	}
	chain.On("GetBlockHash", mock.AnythingOfType("int64")).Return(blockHash, nil)
	chain.On("GetBlockHeader", blockHash).Return(blockHeader, nil)
	var key [16]byte
	emptyFilter, _ := gcs.BuildGCSFilter(0, 0, key, [][]byte{})
	chain.On("GetCFilter", mock.AnythingOfType("chainhash.Hash")).Return(emptyFilter, nil)
	balance, err := scanner.ScanLedger(w)
	assert.NoError(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, uint64(0), balance.Balance)
}

func TestFindBlockHeightFromTime_BinarySearch(t *testing.T) {
	chain, scanner := setupTest()
	targetTime := time.Now()
	blockHash := &chainhash.Hash{}
	chain.On("GetBlockHash", mock.AnythingOfType("int64")).Return(blockHash, nil)
	chain.On("GetBlockHeader", blockHash).Return(&wire.BlockHeader{
		Timestamp: targetTime.Add(-time.Hour),
	}, nil).Once()
	chain.On("GetBlockHeader", blockHash).Return(&wire.BlockHeader{
		Timestamp: targetTime.Add(time.Hour),
	}, nil)
	height, err := scanner.findBlockHeightFromTime(targetTime, 100)
	assert.NoError(t, err)
	assert.True(t, height >= 0)
}

func TestGenerateAllAddresses(t *testing.T) {
	_, scanner := setupTest()
	w := wallet.New(&seed, passphrase, "test")
	w.NextReceiveIndex = 3
	w.NextChangeIndex = 3
	addresses, err := scanner.DeriveAddressSpace(w)
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
	w := wallet.New(&seed, passphrase, "test")
	chain := &MockChainService{}
	scanner := NewBalance(chain)
	addresses, err := scanner.DeriveAddressSpace(w)
	assert.NoError(t, err)
	expectedCount := 21 * 2 // Default 20 addresses + index 0, both receive and change
	assert.Equal(t, expectedCount, len(addresses))
}

func TestAddressesToScripts(t *testing.T) {
	w := wallet.New(&seed, passphrase, "test")
	chain := &MockChainService{}
	scanner := NewBalance(chain)
	addresses, err := scanner.DeriveAddressSpace(w)
	assert.NoError(t, err)
	scripts, err := scanner.addressesToScripts(addresses)
	assert.NoError(t, err)
	assert.Equal(t, len(addresses), len(scripts))
	for _, script := range scripts {
		assert.True(t, len(script) > 0)
	}
}

func TestScanBlockForAddresses_NoMatches(t *testing.T) {
	chain, scanner := setupTest()
	blockHash := &chainhash.Hash{}
	var key [16]byte
	emptyFilter, _ := gcs.BuildGCSFilter(0, 0, key, [][]byte{})
	chain.On("GetCFilter", *blockHash).Return(emptyFilter, nil)
	scripts := [][]byte{
		{0x51, 0x20, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20},
	}
	matches, err := scanner.scanBlockForAddresses(blockHash, scripts)
	assert.NoError(t, err)
	assert.Equal(t, 0, matches)
}

func TestScanBlockForAddresses_GetCFilterError(t *testing.T) {
	chain, scanner := setupTest()
	blockHash := &chainhash.Hash{}
	chain.On("GetCFilter", *blockHash).Return((*gcs.Filter)(nil), assert.AnError)
	scripts := [][]byte{{0x51, 0x20}}
	matches, err := scanner.scanBlockForAddresses(blockHash, scripts)
	assert.Error(t, err)
	assert.Equal(t, 0, matches)
	assert.Contains(t, err.Error(), "failed to get compact filter")
}
