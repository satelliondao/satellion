package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/stdout"
)

type homeModel struct {
	ctx     *AppContext
	cursor  int
	choices []string
}

var (
	MenuSyncBlockchain  = "Sync blockchain"
	MenuCreateNewWallet = "Create new wallet"
	MenuListWallets     = "List wallets"
)

var choices = []string{
	MenuSyncBlockchain,
	MenuCreateNewWallet,
	MenuListWallets,
}

func NewHome(ctx *AppContext) Page {
	return &homeModel{ctx: ctx, choices: choices}
}

func (m *homeModel) Init() tea.Cmd { return nil }

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
				return m, Navigate("create")
			case MenuListWallets:
				return m, Navigate("list")
			// 	return m, nil
			// case "Import wallet from seed":
			// 	m.ctx.Router.ImportWalletFromSeed()
			// 	return m, nil
			// case "Show wallet info":
			// 	m.ctx.Router.ShowWalletInfo()
			// 	return m, nil
			// case "List wallets":
			// 	m.ctx.Router.ListWallets()
			// return m, nil
			case MenuSyncBlockchain:
				return m, Navigate("sync")
			}
		}
	}
	return m, nil
}

func (m *homeModel) View() string {
	title := stdout.Title() + "\n\n"
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
