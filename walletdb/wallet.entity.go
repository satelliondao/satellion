package walletdb

import "github.com/satelliondao/satellion/wallet"

type WalletEntity struct {
	Name             string   `json:"name"`
	Mnemonic         []string `json:"mnemonic"`
	NextChangeIndex  uint32   `json:"next_change_index"`
	NextReceiveIndex uint32   `json:"next_receive_index"`
}

func NewWalletEntity(w *wallet.Wallet) *WalletEntity {
	if w.RootKey == nil {
		panic("root key is nil")
	}
	if w.Mnemonic == nil {
		panic("mnemonic is nil")
	}
	return &WalletEntity{
		Name:             w.Name,
		Mnemonic:         w.Mnemonic.Words,
		NextChangeIndex:  w.NextChangeIndex,
		NextReceiveIndex: w.NextReceiveIndex,
	}
}
