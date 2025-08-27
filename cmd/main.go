package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/ui"
	"github.com/satelliondao/satellion/usecase"
)

func main() {
	r := usecase.NewRouter()
	ctx := ui.NewContext(r)
	pages := map[string]ui.PageFactory{
		config.HomePage:           ui.NewHome,
		config.SyncPage:           ui.NewSync,
		config.CreateWalletPage:   ui.NewCreateWallet,
		config.VerifyMnemonicPage: ui.NewVerifyMnemonic,
		config.ListWalletsPage:    ui.NewListWallets,
	}
	app := ui.NewApp(ctx, pages, config.HomePage)
	_, _ = tea.NewProgram(app, tea.WithAltScreen()).Run()
}
