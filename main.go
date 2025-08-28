package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/router"
	"github.com/satelliondao/satellion/ui/frame"
	"github.com/satelliondao/satellion/ui/home"
	"github.com/satelliondao/satellion/ui/sync"
	"github.com/satelliondao/satellion/ui/wallet_create"
	"github.com/satelliondao/satellion/ui/wallet_list"
	"github.com/satelliondao/satellion/ui/wallet_switch"
)

func main() {
	r := router.NewRouter()
	ctx := frame.NewContext(r)
	pages := map[string]frame.PageFactory{
		config.HomePage:           home.New,
		config.SyncPage:           sync.New,
		config.CreateWalletPage:   wallet_create.New,
		config.VerifyMnemonicPage: wallet_create.NewVerify,
		config.ListWalletsPage:    wallet_list.New,
		config.SwitchWalletPage:   wallet_switch.New,
	}
	app := frame.NewApp(ctx, pages, config.HomePage)
	_, _ = tea.NewProgram(app, tea.WithAltScreen()).Run()
}
