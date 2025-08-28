package home

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/frame"
)

type state struct {
	ctx          *frame.AppContext
	cursor       int
	choices      []string
	activeWallet string
	walletsCount int
}

type menuItem struct{ label, page string }

var menuItems = []menuItem{
	{label: "Receive", page: config.ReceivePage},
	{label: "Send", page: config.SendPage},
	{label: "Syncronize chain", page: config.SyncPage},
	{label: "Create new wallet", page: config.CreateWalletPage},
	{label: "List wallets", page: config.ListWalletsPage},
	{label: "Switch active wallet", page: config.SwitchWalletPage},
}

func New(ctx *frame.AppContext) frame.Page {
	labels := make([]string, len(menuItems))
	for i := range menuItems {
		labels[i] = menuItems[i].label
	}
	return &state{ctx: ctx, choices: labels}
}

func (m *state) Init() tea.Cmd {
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

func (m *state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, frame.Navigate(menuItems[m.cursor].page)
		}
	}
	return m, nil
}

func (m *state) View() string {
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
