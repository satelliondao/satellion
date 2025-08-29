package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/router"
	"github.com/satelliondao/satellion/ui/frame"
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
		config.HomePage:           home.New,
		config.SyncPage:           sync.New,
		config.CreateWalletPage:   wallet_create.New,
		config.VerifyMnemonicPage: wallet_create.NewVerify,
		config.PassphrasePage:     wallet_create.NewPassphrase,
		config.ListWalletsPage:    wallet_list.New,
		config.SwitchWalletPage:   wallet_switch.New,
		config.UnlockWalletPage:   wallet_unlock.New,
		config.ReceivePage:        receive.New,
	}
	startPage := config.UnlockWalletPage

	count, err := ctx.WalletRepo.WalletCount()
	if err == nil && count == 0 {
		startPage = config.CreateWalletPage
	}

	app := frame.NewApp(ctx, pages, startPage)
	_, _ = tea.NewProgram(app, tea.WithAltScreen()).Run()
}
