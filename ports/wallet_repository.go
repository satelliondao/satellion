package ports

import (
	m "github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/wallet"
)

type WalletRepository interface {
	Save(wname string, m *m.Mnemonic) error
	Add(wname string, m *m.Mnemonic) error
	Get(wname string) (*m.Mnemonic, error)
	GetAll() ([]wallet.Wallet, error)
	Delete(wname string) error
	SetDefault(wname string) error
}
