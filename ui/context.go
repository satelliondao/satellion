package ui

import (
	"github.com/satelliondao/satellion/mnemonic"
	usecase "github.com/satelliondao/satellion/router"
)

type AppContext struct {
	Router         *usecase.Router
	TempWalletName string
	TempMnemonic   *mnemonic.Mnemonic
}

func NewContext(router *usecase.Router) *AppContext {
	return &AppContext{Router: router}
}
