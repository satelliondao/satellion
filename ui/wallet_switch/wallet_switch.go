package wallet_switch

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/frame"
	"github.com/satelliondao/satellion/ui/frame/page"
	"github.com/satelliondao/satellion/wallet"
)

type state struct {
	ctx     *frame.AppContext
	wallets []wallet.Wallet
	active  string
	cursor  int
	err     string
}

func New(ctx *frame.AppContext) frame.Page { return &state{ctx: ctx} }

func (m *state) Init() tea.Cmd {
	wallets, err := m.ctx.WalletRepo.GetAll()
	if err != nil {
		m.err = err.Error()
		return nil
	}
	m.wallets = wallets
	name, derr := m.ctx.WalletRepo.GetActiveWalletName()
	if derr == nil {
		m.active = name
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
				m.cursor = len(m.wallets) - 1
			}
		case "down", "j":
			if m.cursor < len(m.wallets)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter":
			if len(m.wallets) == 0 {
				return m, frame.Navigate(page.Home)
			}
			selected := m.wallets[m.cursor].Name
			if err := m.ctx.WalletRepo.SetDefault(selected); err != nil {
				m.err = err.Error()
				return m, nil
			}
			m.ctx.TempPassphrase = ""
			return m, frame.Navigate(page.UnlockWallet)
		}
	}
	return m, nil
}

func (m *state) View() string {
	v := frame.NewViewBuilder()
	v.Line("Switch active wallet")
	if m.err != "" {
		v.Line(m.err)
	}
	if len(m.wallets) == 0 {
		v.Line("No wallets found")
	}
	for i, w := range m.wallets {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		label := w.Name
		if w.Name == m.active {
			label += " (active)"
		}
		v.Line(fmt.Sprintf("%s %s", cursor, label))
	}
	return v.Build()
}
