package ui

import (
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/usecase"
)

type AppContext struct {
	Router         *usecase.Router
	TempWalletName string
	TempMnemonic   *mnemonic.Mnemonic
}

func NewContext(router *usecase.Router) *AppContext {
	return &AppContext{Router: router}
}
