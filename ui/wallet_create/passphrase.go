package wallet_create

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/ui/frame"
)

type passphraseState struct {
	ctx        *frame.AppContext
	passInput  textinput.Model
	confirm    textinput.Model
	err        string
	confirming bool
}

func NewPassphrase(ctx *frame.AppContext) frame.Page {
	m := passphraseState{ctx: ctx}
	in := textinput.New()
	in.Focus()
	in.EchoMode = textinput.EchoPassword
	in.EchoCharacter = '•'
	in.CharLimit = 128
	in.Width = 24
	m.passInput = in
	return m
}

func (m passphraseState) Init() tea.Cmd {
	m.passInput.Focus()
	return textinput.Blink
}

func (m passphraseState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.confirming {
				p := m.passInput.Value()
				if p == "" {
					if err := m.ctx.Router.AddWallet(m.ctx.TempWalletName, *m.ctx.TempMnemonic, ""); err != nil {
						m.err = err.Error()
						return m, nil
					}
					m.ctx.TempWalletName = ""
					m.ctx.TempMnemonic = nil
					return m, frame.Navigate(config.HomePage)
				}
				m.confirm = textinput.New()
				m.confirm.Placeholder = "Confirm passphrase"
				m.confirm.EchoMode = textinput.EchoPassword
				m.confirm.EchoCharacter = '•'
				m.confirm.CharLimit = 128
				m.confirm.Width = 24
				m.confirm.Focus()
				m.confirming = true
				return m, nil
			}
			if m.confirm.Value() != m.passInput.Value() {
				m.err = "Passphrases do not match."
				return m, nil
			}
			if err := m.ctx.Router.AddWallet(m.ctx.TempWalletName, *m.ctx.TempMnemonic, m.passInput.Value()); err != nil {
				m.err = err.Error()
				return m, nil
			}
			m.ctx.TempWalletName = ""
			m.ctx.TempMnemonic = nil
			return m, frame.Navigate(config.HomePage)
		}
	}
	if m.confirming {
		m.confirm, cmd = m.confirm.Update(msg)
		return m, cmd
	}
	m.passInput, cmd = m.passInput.Update(msg)
	return m, cmd
}

func (m passphraseState) View() string {
	v := frame.NewViewBuilder()
	if !m.confirming {
		v.Line("Set an optional passphrase. Leave empty if none.")
		v.Line(m.passInput.View())
		if m.err != "" {
			v.Line(m.err)
		}
		v.Line("Press Enter to continue")
		return v.Build()
	}
	v.Line("Confirm your passphrase.")
	v.Line(m.confirm.View())
	if m.err != "" {
		v.Line(m.err)
	}
	return v.WithHelpText("Press Enter to save").Build()
}
