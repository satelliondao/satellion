package wallet_create

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/ui/page"
	"github.com/satelliondao/satellion/ui/staff"
)

type state struct {
	ctx                *staff.AppContext
	nameInput          textinput.Model
	nameInputCompleted bool
	mnemonic           *mnemonic.Mnemonic
}

func initialState(ctx *staff.AppContext) state {
	i := textinput.New()
	i.Placeholder = "Enter wallet name"
	i.Focus()
	i.CharLimit = 50
	i.Width = 20
	return state{ctx: ctx, nameInput: i}
}

func New(ctx *staff.AppContext) staff.Page {
	return initialState(ctx)
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
				m.nameInputCompleted = true
				m.mnemonic = mnemonic.NewRandom()
			} else {
				if m.ctx != nil {
					m.ctx.TempWalletName = m.nameInput.Value()
					m.ctx.TempMnemonic = m.mnemonic
				}
				return m, staff.Navigate(page.VerifyMnemonic)
			}
		}

		if !m.nameInputCompleted {
			m.nameInput, cmd = m.nameInput.Update(msg)
		}
		return m, cmd
	}

	return m, cmd
}

func (m state) View() string {
	v := staff.NewViewBuilder()
	if m.mnemonic == nil {
		v.Line("Get name for your wallet")
		v.Line(m.nameInput.View())
		return v.Build()
	}

	if m.mnemonic != nil {
		v.Line(fmt.Sprintf("Wallet name: %s", m.nameInput.Value()))
		v.Line(fmt.Sprintf("\n🔑 %s 🔑\n", m.mnemonic.String()))
		v.Line(color.New(color.FgHiRed).Sprintf("Write down your private key and keep it in a safe place."))
		v.Line("You will be asked to verify it in the next step.")
		v.Line("Press enter to continue")
	}

	return v.WithQuitText().Build()
}
