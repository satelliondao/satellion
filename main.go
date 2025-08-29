package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/router"
	"github.com/satelliondao/satellion/ui/frame"
	"github.com/satelliondao/satellion/ui/frame/page"
	"github.com/satelliondao/satellion/ui/home"
	"github.com/satelliondao/satellion/ui/receive"
	"github.com/satelliondao/satellion/ui/sync"
	"github.com/satelliondao/satellion/ui/wallet_create"
	"github.com/satelliondao/satellion/ui/wallet_list"
	"github.com/satelliondao/satellion/ui/wallet_switch"
	"github.com/satelliondao/satellion/ui/wallet_unlock"
)

func main() {
	r := router.NewRouter()
	ctx := frame.NewContext(r)
	pages := map[string]frame.PageFactory{
		page.Home:           home.New,
		page.Sync:           sync.New,
		page.CreateWallet:   wallet_create.New,
		page.VerifyMnemonic: wallet_create.NewVerify,
		page.Passphrase:     wallet_create.NewPassphrase,
		page.ListWallets:    wallet_list.New,
		page.SwitchWallet:   wallet_switch.New,
		page.UnlockWallet:   wallet_unlock.New,
		page.Receive:        receive.New,
	}
	startPage := page.UnlockWallet

	count, err := ctx.WalletRepo.WalletCount()
	if err == nil && count == 0 {
		startPage = page.CreateWallet
	}

	app := frame.NewApp(ctx, pages, startPage)
	_, _ = tea.NewProgram(app, tea.WithAltScreen()).Run()
}
