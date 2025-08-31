package wallet_list

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/router"
	"github.com/satelliondao/satellion/wallet"
)

type state struct {
	ctx     *framework.AppContext
	wallets []wallet.Wallet
}
type errorMsg struct {
	err error
}

func New(ctx *framework.AppContext, params interface{}) framework.Page {
	return &state{ctx: ctx}
}

func (m *state) Init() tea.Cmd {
	wallets, err := m.ctx.WalletRepo.GetAll()
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
			return m, router.Home()
		}
	}
	return m, nil
}

func (m *state) View() string {
	v := framework.NewViewBuilder()

	for i, w := range m.wallets {
		mn, err := m.ctx.WalletRepo.Get(w.Name, "")
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
