package home

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/frame"
	"github.com/satelliondao/satellion/wallet"
)

type state struct {
	ctx     *frame.AppContext
	cursor  int
	choices []string
	w       *wallet.Wallet
	items   []menuItem
}

type menuItem struct{ label, page string }
type errorMsg struct {
	err error
}

var baseMenuItems = []menuItem{
	{label: "Receive", page: config.ReceivePage},
	{label: "Send", page: config.SendPage},
	{label: "Sync chain", page: config.SyncPage},
}

func New(ctx *frame.AppContext) frame.Page {
	m := &state{ctx: ctx}
	m.rebuildMenu()
	return m
}

func (m *state) rebuildMenu() {
	items := make([]menuItem, 0, len(baseMenuItems)+1)
	items = append(items, baseMenuItems...)
	m.items = items
	m.choices = make([]string, len(items))
	for i := range items {
		m.choices[i] = items[i].label
	}
	if m.cursor >= len(m.choices) {
		m.cursor = 0
	}
}

func (m *state) Init() tea.Cmd {
	wallet, err := m.ctx.Router.WalletRepo.GetActiveWallet(m.ctx.TempPassphrase)
	if err != nil {
		return func() tea.Msg { return errorMsg{err: err} }
	}
	m.w = wallet
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
			if len(m.items) == 0 {
				return m, nil
			}
			selected := m.items[m.cursor]
			return m, frame.Navigate(selected.page)
		}
	}
	return m, nil
}

func (m *state) View() string {
	v := frame.NewViewBuilder()
	if m.w != nil {
		v.Line(fmt.Sprintf("Wallet %s\n", color.New(color.Bold).Sprintf("%s", m.w.Name)))
	}
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		v.Line(fmt.Sprintf("%s %s", cursor, choice))
	}
	return v.Build()
}
