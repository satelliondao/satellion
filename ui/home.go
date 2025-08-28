package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/stdout"
)

type homeModel struct {
	ctx          *AppContext
	cursor       int
	choices      []string
	activeWallet string
	walletsCount int
}

var (
	MenuSyncBlockchain  = "Sync blockchain"
	MenuCreateNewWallet = "Create new wallet"
	MenuListWallets     = "List wallets"
	MenuSwitchWallet    = "Switch active wallet"
)

var choices = []string{
	MenuSyncBlockchain,
	MenuCreateNewWallet,
	MenuListWallets,
	MenuSwitchWallet,
}

func NewHome(ctx *AppContext) Page { return &homeModel{ctx: ctx, choices: choices} }

func (m *homeModel) Init() tea.Cmd {
	wallets, err := m.ctx.Router.WalletRepo.GetAll()
	if err == nil {
		m.walletsCount = len(wallets)
		if m.walletsCount > 1 {
			if name, derr := m.ctx.Router.WalletRepo.GetActiveWalletName(); derr == nil {
				m.activeWallet = name
			}
		}
	}
	return nil
}

func (m *homeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		if stdout.ShouldQuit(v) {
			return m, tea.Quit
		}
		switch v.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.choices) - 1
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter":
			switch m.choices[m.cursor] {
			case MenuCreateNewWallet:
				return m, Navigate(config.CreateWalletPage)
			case MenuListWallets:
				return m, Navigate(config.ListWalletsPage)
			case MenuSyncBlockchain:
				return m, Navigate(config.SyncPage)
			case MenuSwitchWallet:
				return m, Navigate(config.SwitchWalletPage)
			}
		}
	}
	return m, nil
}

func (m *homeModel) View() string {
	title := stdout.Title() + "\n\n"
	if m.walletsCount > 1 && m.activeWallet != "" {
		title += fmt.Sprintf("Active wallet: %s\n\n", m.activeWallet)
	}
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		title += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	title += "\nUse ↑/↓ to navigate, Enter to select, Esc to quit\n"
	return title
}
