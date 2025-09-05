package wallet_create

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/router"
)

type state struct {
	ctx                *framework.AppContext
	nameInput          textinput.Model
	nameInputCompleted bool
	mnemonic           *mnemonic.Mnemonic
	err                string
}

func New(ctx *framework.AppContext, params interface{}) framework.Page {
	return &state{ctx: ctx, nameInput: nameInput()}
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
			if m.mnemonic == nil {
				if m.nameInput.Value() == "" {
					m.err = "Wallet name cannot be empty"
					return m, nil
				}
				m.nameInputCompleted = true
				m.mnemonic = mnemonic.NewRandom()
			} else {
				return m, router.VerifyMnemonic(m.nameInput.Value(), m.mnemonic)
			}
		}

		if !m.nameInputCompleted {
			m.nameInput, cmd = m.nameInput.Update(msg)
			if m.err != "" {
				m.err = ""
			}
		}
		return m, cmd
	}

	return m, cmd
}

func (m state) View() string {
	v := framework.View()
	if m.mnemonic == nil {
		v.L("Create new wallet").
			L(m.nameInput.View())
	}

	if m.mnemonic != nil {
		v.L("Wallet name: %s", m.nameInput.Value()).
			L("\n🔑 %s 🔑\n", m.mnemonic.String()).
			L(color.New(color.FgHiRed).Sprintf("Write down your private key and keep it in a safe place")).
			L("You will be asked to verify it in the next step").
			L("Press enter to continue")
	}

	v.Err(m.err)
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
