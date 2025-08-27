package walletdb

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/btcsuite/btcwallet/walletdb"
	bdb "github.com/btcsuite/btcwallet/walletdb"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/wallet"
)

type WalletRepository struct {
	db walletdb.DB
}

func NewWalletRepository(db walletdb.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (s *WalletRepository) bucketName(name string) []byte {
	return []byte("wallet_" + name)
}

func (s *WalletRepository) Save(wname string, m *mnemonic.Mnemonic) error {
	return s.db.Update(func(tx bdb.ReadWriteTx) error {
		bucket := tx.ReadWriteBucket(s.bucketName(wname))
		if bucket == nil {
			b, createErr := tx.CreateTopLevelBucket(s.bucketName(wname))
			if createErr != nil {
				return createErr
			}
			bucket = b
		}
		entity := WalletEntity{Name: wname, Mnemonic: m.String()}
		out, marshalErr := json.Marshal(entity)
		if marshalErr != nil {
			return marshalErr
		}
		return bucket.Put(s.bucketName(wname), out)
	}, func() {})
}

func (s *WalletRepository) Get(wname string) (*mnemonic.Mnemonic, error) {
	var entity WalletEntity
	err := s.db.View(func(tx bdb.ReadTx) error {
		bucketName := s.bucketName(wname)
		bucket := tx.ReadBucket(bucketName)
		if bucket == nil {
			return nil
		}
		raw := bucket.Get(bucketName)
		if len(raw) == 0 {
			legacy := bucket.Get([]byte(wname))
			if len(legacy) == 0 {
				return nil
			}
			entity.Mnemonic = string(legacy)
			return nil
		}
		return json.Unmarshal(raw, &entity)
	}, func() {})
	if err != nil {
		return nil, err
	}
	if entity.Mnemonic == "" {
		return nil, errors.New("wallet not found")
	}
	return mnemonic.New(strings.Split(entity.Mnemonic, " ")), nil
}

func (s *WalletRepository) Add(wname string, m *mnemonic.Mnemonic) error {
	return s.Save(wname, m)
}

func (s *WalletRepository) GetAll() ([]wallet.Wallet, error) {
	var list []wallet.Wallet
	err := s.db.View(func(tx bdb.ReadTx) error {
		return tx.ForEachBucket(func(k []byte) error {
			name := string(k)
			if strings.HasPrefix(name, "wallet_") {
				list = append(list, wallet.Wallet{Name: strings.TrimPrefix(name, "wallet_")})
			}
			return nil
		})
	}, func() {})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *WalletRepository) Delete(wname string) error {
	return s.db.Update(func(tx bdb.ReadWriteTx) error {
		if idx := tx.ReadWriteBucket(walletStoreKey); idx != nil {
			_ = idx.Delete([]byte(wname))
		}
		return tx.DeleteTopLevelBucket(s.bucketName(wname))
	}, func() {})
}

func (s *WalletRepository) SetDefault(wname string) error {
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
