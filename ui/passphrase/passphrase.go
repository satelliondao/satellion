package passphrase

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/router"
)

type State struct {
	ctx        *framework.AppContext
	passInput  textinput.Model
	confirm    textinput.Model
	err        string
	confirming bool
	walletName string
	mnemonic   *mnemonic.Mnemonic
}

func New(ctx *framework.AppContext, p interface{}) framework.Page {
	m := State{ctx: ctx}
	if props, ok := p.(*router.VerifyMnemonicProps); ok {
		m.walletName = props.WalletName
		m.mnemonic = props.Mnemonic
	}
	m.passInput = PassphraseInput("Enter a passphrase (optional)")
	m.confirm = PassphraseInput("Confirm passphrase")
	return m
}

func PassphraseInput(placeholder string) textinput.Model {
	in := textinput.New()
	in.Placeholder = placeholder
	in.EchoMode = textinput.EchoPassword
	in.EchoCharacter = '•'
	in.CharLimit = 128
	in.Width = 24
	in.Focus()
	return in
}

func (m State) Init() tea.Cmd {
	m.passInput.Focus()
	return textinput.Blink
}

func (m State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.confirming {
				m.confirming = true
				return m, nil
			}
			return m.handleConfirmInput()
		}
	}

	if m.confirming {
		m.confirm, cmd = m.confirm.Update(msg)
	} else {
		m.passInput, cmd = m.passInput.Update(msg)
	}

	return m, cmd
}

func (m State) handleConfirmInput() (tea.Model, tea.Cmd) {
	if m.confirm.Value() != m.passInput.Value() {
		m.err = "Passphrases do not match."
		return m, nil
	}
	if err := m.ctx.WalletService.AddWallet(m.walletName, *m.mnemonic, m.passInput.Value()); err != nil {
		m.err = err.Error()
		return m, nil
	}
	m.ctx.Passphrase = m.passInput.Value()
	return m, router.Home()
}

func (m State) View() string {
	v := framework.NewViewBuilder()
	if m.confirming {
		v.Line("Confirm your passphrase:")
		v.Line(m.confirm.View())
	} else {
		v.Line("Enter a passphrase (optional):")
		v.Line(m.passInput.View())
	}
	v.WithErrText(m.err)
	return v.Build()
}
