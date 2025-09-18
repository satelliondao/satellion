package receive

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/router"
	"github.com/satelliondao/satellion/wallet"
)

type state struct {
	ctx     *framework.AppContext
	err     string
	address *wallet.Address
	wallet  *wallet.Wallet
}

type errorMsg struct {
	err string
}

func New(ctx *framework.AppContext, params interface{}) framework.Page {
	s := &state{ctx: ctx}
	return s
}

func (s *state) Init() tea.Cmd {
	w, err := s.ctx.WalletRepo.GetActiveWallet(s.ctx.Passphrase)
	if err != nil || w == nil {
		s.err = fmt.Errorf("wallet not available").Error()
		return nil
	}
	s.wallet = w
	addr, err := w.ReceiveAddress()
	if err != nil {
		s.err = err.Error()
		return nil
	}
	s.address = addr
	return nil
}

func (s *state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	nav := framework.HandleNav(msg, router.Home())
	if nav != nil {
		return s, nav
	}

	switch v := msg.(type) {
	case errorMsg:
		s.err = v.err
		return s, nil
	case tea.KeyMsg:
		if strings.ToLower(v.String()) == "r" {
			return s, s.regenerateAddress()
		}
	}
	return s, nil
}

func (s *state) regenerateAddress() tea.Cmd {
	return func() tea.Msg {
		addr, err := s.wallet.NewReceiveAddress()
		if err != nil {
			return errorMsg{err: err.Error()}
		}
		s.address = addr
		err = s.ctx.WalletRepo.Save(s.wallet)
		if err != nil {
			return errorMsg{err: err.Error()}
		}
		return nil
	}
}

func (s *state) View() string {
	return framework.View().
		L("Address:").
		L(color.New(color.FgGreen).Sprintf(s.address.Address.String())).
		L("Derivation Index: %d", s.address.DeriviationIndex).
		L("").
		Err(s.err).
		Help("R to generate new address").
		QuitHint().
		Build()
}
