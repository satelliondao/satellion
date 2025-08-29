package frame

import (
	"github.com/satelliondao/satellion/mnemonic"
	usecase "github.com/satelliondao/satellion/router"
	"github.com/satelliondao/satellion/walletdb"
)

type AppContext struct {
	Router         *usecase.Router
	WalletRepo     *walletdb.WalletDB
	TempWalletName string
	TempMnemonic   *mnemonic.Mnemonic
	TempPassphrase string
}

func NewContext(router *usecase.Router) *AppContext {
	return &AppContext{Router: router, WalletRepo: router.WalletRepo}
}
