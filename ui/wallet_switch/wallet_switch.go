package wallet_switch

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/router"
	"github.com/satelliondao/satellion/wallet"
)

type state struct {
	ctx     *framework.AppContext
	wallets []wallet.Wallet
	active  string
	cursor  int
	err     string
}

func New(ctx *framework.AppContext, params interface{}) framework.Page {
	return &state{ctx: ctx}
}

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
				return m, router.Home()
			}
			selected := m.wallets[m.cursor].Name
			if err := m.ctx.WalletRepo.SetDefault(selected); err != nil {
				m.err = err.Error()
				return m, nil
			}
			return m, router.UnlockWallet()
		}
	}
	return m, nil
}

func (m *state) View() string {
	v := framework.View().
		L("Switch active wallet")
	if len(m.wallets) == 0 {
		v.L("No wallets found")
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
		v.L("%s %s", cursor, label)
	}
	return v.Err(m.err).
		QuitHint().
		Build()
}
