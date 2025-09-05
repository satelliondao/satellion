package neutrino

import (
	"fmt"
	"log"
	"time"

	"github.com/btcsuite/btcd/btcutil/gcs/builder"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/satelliondao/satellion/ports"
	"github.com/satelliondao/satellion/wallet"
)

const MaxIndexFurtherLookup = 20

type BalanceInfo struct {
	Balance   uint64
	UtxoCount uint64
}

type BalanceService struct {
	chain      ports.Chain
	onProgress func(current, total int64, percent float64)
}

func NewBalance(chain ports.Chain) *BalanceService {
	return &BalanceService{chain: chain}
}

func (s *BalanceService) SetProgressCallback(callback func(current, total int64, percent float64)) {
	s.onProgress = callback
}

func (s *BalanceService) ScanLedger(wallet *wallet.Wallet) (*BalanceInfo, error) {
	if wallet.CreatedAt.IsZero() {
		return nil, fmt.Errorf("wallet creation time not set")
	}
	block, err := s.chain.BestBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to get best block: %w", err)
	}
	startHeight, err := s.findBlockHeightFromTime(wallet.CreatedAt, block.Height)
	if err != nil {
		return nil, fmt.Errorf("failed to find start height: %w", err)
	}
	blockCount := int64(block.Height) - startHeight + 1
	addresses, err := s.DeriveAddressSpace(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to generate addresses: %w", err)
	}
	scripts, err := s.addressesToScripts(addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to convert addresses to scripts: %w", err)
	}
	var totalBalance uint64
	var totalUtxos uint64
	processed := int64(0)
	for height := startHeight; height <= int64(block.Height); height++ {
		res, err := s.scanBlock(height, scripts, processed, blockCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan block %d: %w", height, err)
		}
		totalBalance += res.Balance
		totalUtxos += res.Utxos
		processed++
	}
	return &BalanceInfo{
		Balance:   totalBalance,
		UtxoCount: totalUtxos,
	}, nil
}

type ScanBlockResult struct {
	Balance uint64
	Utxos   uint64
}

func (s *BalanceService) scanBlock(height int64, scripts [][]byte, processed int64, blockCount int64) (*ScanBlockResult, error) {
	blockHash, err := s.chain.GetBlockHash(height)
	res := &ScanBlockResult{
		Balance: 0,
		Utxos:   0,
	}
	if err != nil {
		log.Printf("Warning: failed to get block hash for height %d: %v", height, err)
		return res, err
	}
	matches, err := s.scanBlockForAddresses(blockHash, scripts)
	if err != nil {
		log.Printf("Warning: failed to scan block %d: %v", height, err)
		return res, err
	}
	if matches > 0 {
		res.Utxos = uint64(matches)
		res.Balance = uint64(matches * 1000)
	}
	processed++
	if processed%1000 == 0 || processed == blockCount {
		progress := float64(processed) / float64(blockCount) * 100
		if s.onProgress != nil {
			s.onProgress(processed, blockCount, progress)
		}
	}
	return res, nil
}

func (s *BalanceService) findBlockHeightFromTime(createdAt time.Time, bestHeight int32) (int64, error) {
	var left, right int64 = 0, int64(bestHeight)
	for left < right {
		mid := (left + right) / 2
		blockHash, err := s.chain.GetBlockHash(mid)
		if err != nil {
			return 0, err
		}
		header, err := s.chain.GetBlockHeader(blockHash)
		if err != nil {
			return 0, err
		}
		if header.Timestamp.Before(createdAt) {
			left = mid + 1
		} else {
			right = mid
		}
	}
	if left > 0 {
		left--
	}
	return left, nil
}

func (s *BalanceService) DeriveAddressSpace(w *wallet.Wallet) ([]*wallet.Address, error) {
	var addresses []*wallet.Address
	maxIndex := w.NextReceiveIndex
	if w.NextChangeIndex > maxIndex {
		maxIndex = w.NextChangeIndex
	}
	if maxIndex == 0 {
		maxIndex = MaxIndexFurtherLookup
	}
	for i := uint32(0); i <= maxIndex; i++ {
		receiveAddr, err := w.DeriveTaprootAddress(0, i)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, receiveAddr)
		changeAddr, err := w.DeriveTaprootAddress(1, i)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, changeAddr)
	}
	return addresses, nil
}

func (s *BalanceService) addressesToScripts(addresses []*wallet.Address) ([][]byte, error) {
	var scripts [][]byte
	for _, addr := range addresses {
		script, err := txscript.PayToAddrScript(addr.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to create script for address %s: %w", addr.Address.String(), err)
		}
		scripts = append(scripts, script)
	}
	return scripts, nil
}

func (s *BalanceService) scanBlockForAddresses(blockHash *chainhash.Hash, scripts [][]byte) (int, error) {
	filter, err := s.chain.GetCFilter(*blockHash)
	if err != nil {
		return 0, fmt.Errorf("failed to get compact filter: %w", err)
	}

	key := builder.DeriveKey(blockHash)
	matches := 0
	for _, script := range scripts {
		match, err := filter.Match(key, script)
		if err != nil {
			return 0, fmt.Errorf("failed to execute match on compact filter: %w", err)
		}
		if match {
			matches++
		}
	}

	return matches, nil
}
