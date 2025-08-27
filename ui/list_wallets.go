package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/cfg"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/wallet"
)

type listWalletsModel struct {
	ctx     *AppContext
	wallets []wallet.Wallet
}

func NewListWallets(ctx *AppContext) Page {
	return &listWalletsModel{ctx: ctx}
}

func (m *listWalletsModel) Init() tea.Cmd {
	wallets, err := m.ctx.Router.WalletRepo.GetAll()
	if err != nil {
		return func() tea.Msg { return errorMsg{err: err} }
	}
	m.wallets = wallets
	return nil
}

func (m *listWalletsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		if stdout.ShouldQuit(v) {
			return m, tea.Quit
		}
		if v.String() == "enter" {
			return m, Navigate(cfg.HomePage)
		}
		if v.String() == "up" {
			return m, Navigate(cfg.HomePage)
		}
		if v.String() == "down" {
			return m, Navigate(cfg.HomePage)
		}
	}
	return m, nil
}

func (m *listWalletsModel) View() string {
	view := "Wallet List\n\n"
	for i, w := range m.wallets {
		mn, err := m.ctx.Router.WalletRepo.Get(w.Name)
		mnemonicText := "<not found>"
		if err == nil && mn != nil {
			mnemonicText = mn.String()
		}
		view += fmt.Sprintf("%d. %s\n   %s\n", i+1, w.Name, mnemonicText)
	}
	return view
}
