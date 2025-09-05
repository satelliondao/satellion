package wallet_import

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/router"
)

type state struct {
	ctx                 *framework.AppContext
	nameInput           textinput.Model
	mnemonicInput       textinput.Model
	passphraseInput     textinput.Model
	nameCompleted       bool
	mnemonicCompleted   bool
	passphraseCompleted bool
	err                 string
}

func New(ctx *framework.AppContext, params interface{}) framework.Page {
	return &state{
		ctx:             ctx,
		nameInput:       nameInput(),
		mnemonicInput:   mnemonicInput(),
		passphraseInput: passphraseInput("Enter a passphrase (optional)"),
	}
}

func (m state) Init() tea.Cmd {
	return textinput.Blink
}

func (m state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			return m.handleEnter()
		}

		if !m.nameCompleted {
			m.nameInput, cmd = m.nameInput.Update(msg)
		} else if !m.mnemonicCompleted {
			m.mnemonicInput, cmd = m.mnemonicInput.Update(msg)
		} else if !m.passphraseCompleted {
			m.passphraseInput, cmd = m.passphraseInput.Update(msg)
		}
		return m, cmd
	}

	return m, cmd
}

func (m state) handleEnter() (tea.Model, tea.Cmd) {
	if !m.nameCompleted {
		if m.nameInput.Value() == "" {
			m.err = "Wallet name cannot be empty"
			return m, nil
		}
		m.nameCompleted = true
		m.mnemonicInput.Focus()
		return m, nil
	}

	if !m.mnemonicCompleted {
		if m.mnemonicInput.Value() == "" {
			m.err = "Mnemonic cannot be empty"
			return m, nil
		}
		m.mnemonicCompleted = true
		m.passphraseInput.Focus()
		return m, nil
	}

	if !m.passphraseCompleted {
		if err := m.ctx.WalletService.ImportWallet(m.nameInput.Value(), m.mnemonicInput.Value(), m.passphraseInput.Value()); err != nil {
			m.err = err.Error()
			return m, nil
		}

		m.ctx.Passphrase = m.passphraseInput.Value()
		return m, router.Home()
	}

	return m, nil
}

func (m state) View() string {
	v := framework.View()

	if !m.nameCompleted {
		v.L("Import existing wallet").
			L("Enter wallet name:").
			L(m.nameInput.View())
	} else if !m.mnemonicCompleted {
		v.L("Import wallet").
			L("Wallet name: %s", m.nameInput.Value()).
			L("Enter your 12-word mnemonic phrase:").
			L(m.mnemonicInput.View())
	} else if !m.passphraseCompleted {
		v.L("Import wallet").
			L("Wallet name: %s", m.nameInput.Value()).
			L("Mnemonic: %s", m.mnemonicInput.Value()).
			L("Enter passphrase:").
			L(m.passphraseInput.View())
	}

	if m.err != "" {
		v.L("Error: %s", m.err)
	}

	return v.QuitHint().Build()
}

func nameInput() textinput.Model {
	i := textinput.New()
	i.Placeholder = "Enter wallet name"
	i.Focus()
	i.CharLimit = 50
	i.Width = 20
	return i
}

func mnemonicInput() textinput.Model {
	i := textinput.New()
	i.Placeholder = "Enter 12-word mnemonic phrase"
	i.CharLimit = 200
	i.Width = 50
	return i
}

func passphraseInput(placeholder string) textinput.Model {
	i := textinput.New()
	i.Placeholder = placeholder
	i.EchoMode = textinput.EchoPassword
	i.EchoCharacter = '•'
	i.CharLimit = 128
	i.Width = 24
	return i
}
