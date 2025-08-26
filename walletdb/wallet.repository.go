package walletdb

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/btcsuite/btcwallet/walletdb"
	bdb "github.com/btcsuite/btcwallet/walletdb"
	"github.com/satelliondao/satellion/mnemonic"
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
