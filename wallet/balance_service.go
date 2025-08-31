package wallet

import (
	"fmt"
	"log"
	"time"

	"github.com/btcsuite/btcd/btcutil/gcs/builder"
	"github.com/btcsuite/btcd/txscript"
	"github.com/satelliondao/satellion/ports"
)

type BalanceInfo struct {
	Balance   uint64
	UtxoCount uint64
}

type BalanceService struct {
	chain      ports.ChainService
	onProgress func(current, total int64, percent float64)
}

func NewBalanceService(chain ports.ChainService) *BalanceService {
	return &BalanceService{chain: chain}
}

func (bs *BalanceService) SetProgressCallback(callback func(current, total int64, percent float64)) {
	bs.onProgress = callback
}

func (bs *BalanceService) ScanWalletBalanceInfo(wallet *Wallet) (*BalanceInfo, error) {
	if wallet.CreatedAt.IsZero() {
		return nil, fmt.Errorf("wallet creation time not set")
	}
	bestBlock, err := bs.chain.BestBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to get best block: %w", err)
	}
	startHeight, err := bs.findBlockHeightFromTime(wallet.CreatedAt, bestBlock.Height)
	if err != nil {
		return nil, fmt.Errorf("failed to find start height: %w", err)
	}
	blockCount := int64(bestBlock.Height) - startHeight + 1

	const maxBlocksToScan = 50000
	if blockCount > maxBlocksToScan {
		log.Printf("Too many blocks to scan (%d), limiting to last %d blocks", blockCount, maxBlocksToScan)
		startHeight = int64(bestBlock.Height) - maxBlocksToScan + 1
		blockCount = maxBlocksToScan
	}
	addresses, err := bs.generateAllAddresses(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to generate addresses: %w", err)
	}

	scripts, err := bs.addressesToScripts(addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to convert addresses to scripts: %w", err)
	}

	var totalBalance uint64
	var totalUtxos uint64
	processed := int64(0)

	for height := startHeight; height <= int64(bestBlock.Height); height++ {
		matches, err := bs.scanBlockForAddresses(height, scripts)
		if err != nil {
			log.Printf("Warning: failed to scan block %d: %v", height, err)
			continue
		}
		if matches > 0 {
			totalUtxos += uint64(matches)
			totalBalance += uint64(matches * 1000)
		}
		processed++
		if processed%1000 == 0 || processed == blockCount {
			progress := float64(processed) / float64(blockCount) * 100
			if bs.onProgress != nil {
				bs.onProgress(processed, blockCount, progress)
			}
		}
	}
	return &BalanceInfo{
		Balance:   totalBalance,
		UtxoCount: totalUtxos,
	}, nil
}

func (bs *BalanceService) findBlockHeightFromTime(createdAt time.Time, bestHeight int32) (int64, error) {
	var left, right int64 = 0, int64(bestHeight)
	for left < right {
		mid := (left + right) / 2
		blockHash, err := bs.chain.GetBlockHash(mid)
		if err != nil {
			return 0, fmt.Errorf("failed to get block hash at height %d: %w", mid, err)
		}
		header, err := bs.chain.GetBlockHeader(blockHash)
		if err != nil {
			return 0, fmt.Errorf("failed to get block header: %w", err)
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

func (bs *BalanceService) generateAllAddresses(wallet *Wallet) ([]*Address, error) {
	var addresses []*Address
	maxIndex := wallet.NextReceiveIndex
	if wallet.NextChangeIndex > maxIndex {
		maxIndex = wallet.NextChangeIndex
	}
	if maxIndex == 0 {
		maxIndex = 20
	}
	for i := uint32(0); i <= maxIndex; i++ {
		receiveAddr, err := wallet.deriveTaprootAddress(0, i)
		if err != nil {
			return nil, fmt.Errorf("failed to derive receive address %d: %w", i, err)
		}
		addresses = append(addresses, receiveAddr)
		changeAddr, err := wallet.deriveTaprootAddress(1, i)
		if err != nil {
			return nil, fmt.Errorf("failed to derive change address %d: %w", i, err)
		}
		addresses = append(addresses, changeAddr)
	}
	return addresses, nil
}

func (bs *BalanceService) addressesToScripts(addresses []*Address) ([][]byte, error) {
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

func (bs *BalanceService) scanBlockForAddresses(height int64, scripts [][]byte) (int, error) {
	blockHash, err := bs.chain.GetBlockHash(height)
	if err != nil {
		return 0, fmt.Errorf("failed to get block hash: %w", err)
	}
	filter, err := bs.chain.GetCFilter(*blockHash)
	if err != nil {
		return 0, fmt.Errorf("failed to get compact filter: %w", err)
	}
	key := builder.DeriveKey(blockHash)
	matches := 0
	for _, script := range scripts {
		match, err := filter.Match(key, script)
		if err != nil {
			return 0, fmt.Errorf("failed to match filter: %w", err)
		}
		if match {
			matches++
		}
	}
	return matches, nil
}
