package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/wallet"
	"github.com/satelliondao/satellion/walletdb"
)

type WalletService struct {
	walletRepo *walletdb.WalletDB
}

func NewWalletService(walletRepo *walletdb.WalletDB) *WalletService {
	return &WalletService{walletRepo: walletRepo}
}

func (s *WalletService) AddWallet(name string, m mnemonic.Mnemonic, passphrase string) error {
	if name == "" {
		return fmt.Errorf("invalid wallet data")
	}
	model := wallet.New(&m, passphrase, "")
	model.Name = name
	model.CreatedAt = time.Now()

	err := s.walletRepo.Save(model)
	if err != nil {
		return err
	}
	return s.walletRepo.SetDefault(name)
}

func (s *WalletService) Unlock(passphrase string) error {
	w, err := s.walletRepo.GetActiveWallet(passphrase)
	if err != nil {
		return err
	}
	if w == nil || w.Mnemonic == nil {
		return fmt.Errorf("no active wallet")
	}
	hashSeed := sha256.Sum256(w.Mnemonic.Seed(passphrase))
	unlockKey := hex.EncodeToString(hashSeed[:])

	if w.Lock != unlockKey {
		return fmt.Errorf("invalid passphrase")
	}
	return nil
}

func (s *WalletService) WalletRepo() *walletdb.WalletDB {
	return s.walletRepo
}
