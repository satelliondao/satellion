package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/cfg"
	"github.com/satelliondao/satellion/ui"
	"github.com/satelliondao/satellion/usecase"
)

func main() {
	r := usecase.NewRouter()
	ctx := ui.NewContext(r)
	pages := map[string]ui.PageFactory{
		cfg.HomePage:         ui.NewHome,
		cfg.SyncPage:         ui.NewSync,
		cfg.CreateWalletPage: ui.NewCreateWallet,
	}
	app := ui.NewApp(ctx, pages, cfg.HomePage)
	_, _ = tea.NewProgram(app).Run()
}
