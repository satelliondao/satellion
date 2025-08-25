package ports

import m "github.com/satelliondao/satellion/mnemonic"

type WalletRepo interface {
	Save(wname string, m *m.Mnemonic) error
}
