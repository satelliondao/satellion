package ports

import (
	"github.com/btcsuite/btcd/btcutil/gcs"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/neutrino/headerfs"
)

// ChainService defines the interface for blockchain operations
type ChainService interface {
	BestBlock() (*headerfs.BlockStamp, error)
	GetBlockHash(height int64) (*chainhash.Hash, error)
	GetBlockHeader(hash *chainhash.Hash) (*wire.BlockHeader, error)
	GetCFilter(hash chainhash.Hash) (*gcs.Filter, error)
}
