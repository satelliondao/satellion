package walletdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcwallet/walletdb"
	bdb "github.com/btcsuite/btcwallet/walletdb"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/wallet"
)

const ActiveWalletKey = "active_wallet"

type WalletDB struct {
	db walletdb.DB
}

func New(db walletdb.DB) *WalletDB {
	return &WalletDB{db: db}
}

func (s *WalletDB) getKey(wname string) []byte {
	return []byte("wallet_" + wname)
}

func (s *WalletDB) Add(w *wallet.Wallet) error {
	return s.Save(w)
}

func (s *WalletDB) Save(w *wallet.Wallet) error {
	key := s.getKey(w.Name)
	return s.db.Update(func(tx bdb.ReadWriteTx) error {
		bucket := tx.ReadWriteBucket(key)
		if bucket == nil {
			b, createErr := tx.CreateTopLevelBucket(key)
			if createErr != nil {
				return createErr
			}
			bucket = b
		}
		out, marshalErr := json.Marshal(NewWalletEntity(w))
		if marshalErr != nil {
			return marshalErr
		}
		return bucket.Put(key, out)
	}, func() {})
}

func (s *WalletDB) Get(wname string, passphrase string) (*wallet.Wallet, error) {
	var entity WalletEntity
	e := errors.New("wallet not found")
	err := s.db.View(func(tx bdb.ReadTx) error {
		key := s.getKey(wname)
		bucket := tx.ReadBucket(key)
		if bucket == nil {
			return e
		}
		raw := bucket.Get(key)
		if len(raw) == 0 {
			return e
		}
		return json.Unmarshal(raw, &entity)
	}, func() {})
	if err != nil {
		return nil, err
	}
	return s.toModel(entity, passphrase), nil
}

func (s *WalletDB) WalletCount() (int, error) {
	var count int
	err := s.db.View(func(tx bdb.ReadTx) error {
		return tx.ForEachBucket(func(k []byte) error {
			if strings.HasPrefix(string(k), "wallet_") {
				count++
			}
			return nil
		})
	}, func() {})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *WalletDB) GetAll() ([]wallet.Wallet, error) {
	var list []wallet.Wallet

	err := s.db.View(func(tx bdb.ReadTx) error {
		return tx.ForEachBucket(func(k []byte) error {
			key := string(k)
			if strings.HasPrefix(key, "wallet_") {
				raw := tx.ReadBucket(k).Get(k)
				if len(raw) == 0 {
					return nil
				}

				entity := WalletEntity{}
				if err := json.Unmarshal(raw, &entity); err != nil {
					fmt.Println("failed to unmarshal wallet: ", err)
					return nil
				}
				list = append(list, *s.toModel(entity, ""))
			}
			return nil
		})
	}, func() {})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *WalletDB) Delete(wname string) error {
	key := s.getKey(wname)
	return s.db.Update(func(tx bdb.ReadWriteTx) error {
		if idx := tx.ReadWriteBucket(key); idx != nil {
			_ = idx.Delete(key)
		}
		return tx.DeleteTopLevelBucket(key)
	}, func() {})
}

func (s *WalletDB) SetDefault(wname string) error {
	return s.db.Update(func(tx bdb.ReadWriteTx) error {
		idx := tx.ReadWriteBucket(walletStoreKey)
		if idx == nil {
			b, err := tx.CreateTopLevelBucket(walletStoreKey)
			if err != nil {
				return err
			}
			idx = b
		}
		return idx.Put([]byte(ActiveWalletKey), []byte(wname))
	}, func() {})
}

func (s *WalletDB) GetActiveWallet(passphrase string) (*wallet.Wallet, error) {
	walletName, err := s.GetActiveWalletName()
	e := errors.New("default wallet not set")
	if err != nil {
		return nil, err
	}

	var entity WalletEntity
	key := s.getKey(walletName)
	err = s.db.View(func(tx bdb.ReadTx) error {
		bucket := tx.ReadBucket(key)
		if bucket == nil {
			return e
		}
		raw := bucket.Get(key)
		if len(raw) == 0 {
			return e
		}
		return json.Unmarshal(raw, &entity)
	}, func() {})
	if err != nil {
		return nil, err
	}
	return s.toModel(entity, passphrase), nil
}

func (s *WalletDB) GetActiveWalletName() (string, error) {
	var walletName string
	e := errors.New("default wallet not set")

	err := s.db.View(func(tx bdb.ReadTx) error {
		idx := tx.ReadBucket(walletStoreKey)
		if idx == nil {
			return e
		}
		raw := idx.Get([]byte(ActiveWalletKey))
		if len(raw) == 0 {
			return e
		}
		walletName = string(raw)
		return nil
	}, func() {})
	if err != nil {
		return "", err
	}
	return walletName, nil
}

func (s *WalletDB) toModel(w WalletEntity, passphrase string) *wallet.Wallet {
	mnemonic := mnemonic.New(w.Mnemonic)
	return wallet.New(&mnemonic, passphrase, w.Name, w.NextChangeIndex, w.NextReceiveIndex, w.Lock)
}
