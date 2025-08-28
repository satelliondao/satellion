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

type WalletDB struct {
	db walletdb.DB
}

func New(db walletdb.DB) *WalletDB {
	return &WalletDB{db: db}
}

func (s *WalletDB) bucketName(name string) []byte {
	return []byte("wallet_" + name)
}

func (s *WalletDB) Add(w *wallet.Wallet) error {
	return s.Save(w)
}

func (s *WalletDB) Save(w *wallet.Wallet) error {
	return s.db.Update(func(tx bdb.ReadWriteTx) error {
		bucket := tx.ReadWriteBucket(s.bucketName(w.Name))
		if bucket == nil {
			b, createErr := tx.CreateTopLevelBucket(s.bucketName(w.Name))
			if createErr != nil {
				return createErr
			}
			bucket = b
		}
		out, marshalErr := json.Marshal(NewWalletEntity(w))
		if marshalErr != nil {
			return marshalErr
		}
		return bucket.Put(s.bucketName(w.Name), out)
	}, func() {})
}

func (s *WalletDB) Get(wname string) (*wallet.Wallet, error) {
	var entity WalletEntity
	err := s.db.View(func(tx bdb.ReadTx) error {
		bucketName := s.bucketName(wname)
		bucket := tx.ReadBucket(bucketName)
		if bucket == nil {
			return errors.New("wallet not found")
		}
		raw := bucket.Get(bucketName)
		if len(raw) == 0 {
			legacy := bucket.Get([]byte(wname))
			if len(legacy) == 0 {
				return errors.New("wallet not found")
			}
			return nil
		}
		return json.Unmarshal(raw, &entity)
	}, func() {})
	if err != nil {
		return nil, err
	}
	return s.toModel(entity), nil
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
				list = append(list, *s.toModel(entity))
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
	key := s.bucketName(wname)
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
		return idx.Put([]byte("__default__"), []byte(wname))
	}, func() {})
}

func (s *WalletDB) toModel(entity WalletEntity) *wallet.Wallet {
	mnemonic := mnemonic.New(entity.Mnemonic)
	return wallet.New(&mnemonic, entity.Name, entity.NextIndex)
}
