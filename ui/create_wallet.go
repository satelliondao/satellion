package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/cfg"
	"github.com/satelliondao/satellion/mnemonic"
)

type model struct {
	ctx                *AppContext
	nameInput          textinput.Model
	nameInputCompleted bool
	mnemonic           *mnemonic.Mnemonic
}

func initialModel(ctx *AppContext) model {
	i := textinput.New()
	i.Placeholder = "Enter wallet name"
	i.Focus()
	i.CharLimit = 50
	i.Width = 20
	return model{ctx: ctx, nameInput: i}
}

func NewCreateWallet(ctx *AppContext) Page {
	return initialModel(ctx)
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.mnemonic == nil {
				m.nameInputCompleted = true
				m.mnemonic = mnemonic.NewRandom()
			} else {
				if m.ctx != nil {
					m.ctx.TempWalletName = m.nameInput.Value()
					m.ctx.TempMnemonic = m.mnemonic
				}
				return m, Navigate(cfg.VerifyMnemonicPage)
			}
		}

		if !m.nameInputCompleted {
			m.nameInput, cmd = m.nameInput.Update(msg)
		}
		return m, cmd
	}

	return m, cmd
}

func (m model) View() string {
	if m.mnemonic == nil {
		return fmt.Sprintf("Get name for your wallet\n\n%s", m.nameInput.View())
	}
	if m.mnemonic != nil {
		return fmt.Sprintf(
			"Wallet name: %s\n\n%s\n\n%s\n%s\n%s\n",
			m.nameInput.Value(),
			color.New(color.FgHiYellow).Sprintf("🔑 %s", m.mnemonic.String()),
			"Please write down your seed phrase and keep it in a safe place.",
			"You will be asked to verify it in the next step.",
			"Press enter to continue",
		)
	}

	return ""
}
