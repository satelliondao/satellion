package neutrino

import (
	"github.com/btcsuite/btcd/btcutil/gcs"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/satelliondao/satellion/ports"
	"github.com/stretchr/testify/mock"
)

type MockChainService struct {
	mock.Mock
}

// Compile-time check to ensure MockChainService implements ports.ChainService
var _ ports.Chain = (*MockChainService)(nil)

func (m *MockChainService) BestBlock() (*ports.BlockInfo, error) {
	args := m.Called()
	return args.Get(0).(*ports.BlockInfo), args.Error(1)
}

func (m *MockChainService) GetBlockHash(height int64) (*chainhash.Hash, error) {
	args := m.Called(height)
	return args.Get(0).(*chainhash.Hash), args.Error(1)
}

func (m *MockChainService) GetBlockHeader(hash *chainhash.Hash) (*wire.BlockHeader, error) {
	args := m.Called(hash)
	return args.Get(0).(*wire.BlockHeader), args.Error(1)
}

func (m *MockChainService) GetCFilter(hash chainhash.Hash) (*gcs.Filter, error) {
	args := m.Called(hash)
	return args.Get(0).(*gcs.Filter), args.Error(1)
}
