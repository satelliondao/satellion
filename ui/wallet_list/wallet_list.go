package wallet_list

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/page"
	"github.com/satelliondao/satellion/ui/staff"
	"github.com/satelliondao/satellion/wallet"
)

type state struct {
	ctx     *staff.AppContext
	wallets []wallet.Wallet
}
type errorMsg struct {
	err error
}

func New(ctx *staff.AppContext) staff.Page {
	return &state{ctx: ctx}
}

func (m *state) Init() tea.Cmd {
	wallets, err := m.ctx.Router.WalletRepo.GetAll()
	if err != nil {
		return func() tea.Msg { return errorMsg{err: err} }
	}
	m.wallets = wallets
	return nil
}

func (m *state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		if stdout.ShouldQuit(v) {
			return m, tea.Quit
		}
		if v.String() == "enter" {
			return m, staff.Navigate(page.Home)
		}
	}
	return m, nil
}

func (m *state) View() string {
	v := staff.NewViewBuilder()

	for i, w := range m.wallets {
		mn, err := m.ctx.Router.WalletRepo.Get(w.Name, "")
		mnemonicText := "<not found>"
		if err == nil {
			mnemonicText = mn.Mnemonic.String()
		}
		v.Line(fmt.Sprintf("%d. %s\n   %s\n", i+1, w.Name, mnemonicText))
	}

	if len(m.wallets) == 0 {
		v.Line("No wallets found")
	}

	v.WithQuitText()
	return v.Build()
}
