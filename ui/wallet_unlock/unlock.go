package wallet_unlock

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/passphrase"
	"github.com/satelliondao/satellion/ui/router"
)

type state struct {
	ctx      *framework.AppContext
	input    textinput.Model
	selector *framework.ChoiceSelector
	err      string
}

var choices = []framework.Choice{
	{Label: "Unlock", Value: "unlock"},
	{Label: "Switch wallet", Value: "switch"},
	{Label: "Create new wallet", Value: "create"},
}

func New(ctx *framework.AppContext, params interface{}) framework.Page {
	return &state{
		ctx:      ctx,
		input:    passphrase.PassphraseInput("Enter your passphrase"),
		selector: framework.NewChoiceSelector(choices),
	}
}

func (m *state) Init() tea.Cmd {
	m.input.Focus()
	return textinput.Blink
}

func (m *state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		if v.Type == tea.KeyCtrlC || v.Type == tea.KeyEsc {
			return m, tea.Quit
		}
		if v.Type == tea.KeyCtrlS && v.String() == "ctrl+s" {
			return m, router.SwitchWallet()
		}
		if v.String() == "up" || v.String() == "down" || v.String() == "k" || v.String() == "j" || v.String() == "enter" {
			choiceResult := m.selector.Update(msg)
			if choiceResult.Consumed {
				if choiceResult.Action == framework.ActionSelection && choiceResult.Selected != nil {
					switch choiceResult.Selected.Value {
					case "unlock":
						pass := m.input.Value()
						if err := m.ctx.WalletService.Unlock(pass); err != nil {
							m.input.SetValue("")
							m.err = err.Error()
							return m, nil
						}
						m.ctx.Passphrase = pass
						return m, router.Home()
					case "switch":
						m.input.SetValue("")
						return m, router.SwitchWallet()
					case "create":
						return m, router.CreateWallet()
					}
				}
				return m, nil
			}
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *state) View() string {
	v := framework.View()
	name, err := m.ctx.WalletRepo.GetActiveWalletName()
	if err != nil {
		v.L("No active wallet found\n")
	} else {
		v.L("Enter passphrase to unlock wallet %s\n", color.New(color.Bold).Sprintf("%s", name))
	}
	v.L(m.input.View())
	v.Err(m.err)
	v.L("")
	v.L(m.selector.RenderWithPrefix("▶"))
	return v.Build()
}
