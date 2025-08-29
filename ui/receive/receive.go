package receive

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/frame"
)

type addressGeneratedMsg struct {
	address string
	index   uint32
	err     error
}

type state struct {
	ctx            *frame.AppContext
	address        string
	index          uint32
	derivationPath string
	err            string
	showAddr       bool
}

func New(ctx *frame.AppContext) frame.Page {
	return &state{
		ctx: ctx,
	}
}

func (s *state) Init() tea.Cmd {
	return s.generateAddress()
}

func (s *state) generateAddress() tea.Cmd {
	return func() tea.Msg {
		w, err := s.ctx.Router.WalletRepo.GetActiveWallet(s.ctx.TempPassphrase)
		if err != nil || w == nil {
			return addressGeneratedMsg{err: fmt.Errorf("wallet not available")}
		}

		addr, derr := w.DeriveReceiveAddress()
		if derr != nil {
			return addressGeneratedMsg{err: derr}
		}
		return addressGeneratedMsg{
			address: addr.Address,
			index:   addr.DeriviationIndex,
		}
	}
}

func (s *state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case addressGeneratedMsg:
		if v.err != nil {
			s.err = v.err.Error()
			return s, nil
		}
		s.address = v.address
		s.index = v.index
		s.showAddr = true
		s.err = ""
		return s, nil

	case tea.KeyMsg:
		if stdout.ShouldQuit(v) {
			return s, frame.Navigate(config.HomePage)
		}

		if v.String() == "r" {
			return s, s.generateAddress()
		}
	}

	return s, nil
}

func (s *state) View() string {
	v := frame.NewViewBuilder()
	v.Line("Receive Bitcoin")
	v.Line("")

	if s.showAddr {
		v.Line(color.New(color.FgGreen).Sprintf("Address: %s", s.address))
		v.Line(color.New(color.FgCyan).Sprintf("Derivation Path: %d", s.index))
		v.Line("")
		v.Line("Press Enter or 'r' to generate a new anonymous address")
	} else {
		v.Line("Generating address...")
	}

	v.WithErrText(s.err)
	return v.Build()
}
