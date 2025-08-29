package receive

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/frame"
	"github.com/satelliondao/satellion/ui/frame/page"
	"github.com/satelliondao/satellion/wallet"
)

type state struct {
	ctx     *frame.AppContext
	err     string
	address *wallet.Address
	wallet  *wallet.Wallet
}

type errorMsg struct {
	err string
}

func New(ctx *frame.AppContext) frame.Page {
	return &state{
		ctx: ctx,
	}
}

func (s *state) Init() tea.Cmd {
	w, err := s.ctx.Router.WalletRepo.GetActiveWallet(s.ctx.TempPassphrase)
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
	switch v := msg.(type) {
	case errorMsg:
		s.err = v.err
		return s, nil
	case tea.KeyMsg:
		if stdout.ShouldQuit(v) {
			return s, tea.Quit
		}
		if v.Type == tea.KeyEsc {
			return s, frame.Navigate(page.Home)
		}
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
	v := frame.NewViewBuilder()
	v.Line(color.New(color.FgGreen).Sprintf("Address: %s", s.address.Address))
	v.Line(color.New(color.FgCyan).Sprintf("Derivation Path: %d", s.address.DeriviationIndex))
	v.Line("")
	v.WithHelpText("R to generate new address")
	v.WithHelpText("Esc or Ctrl+C to go back")
	v.WithErrText(s.err)
	return v.Build()
}
