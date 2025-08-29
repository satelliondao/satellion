package wallet_unlock

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/ui/frame"
)

type state struct {
	ctx      *frame.AppContext
	input    textinput.Model
	selector *frame.ChoiceSelector
	err      string
}

func New(ctx *frame.AppContext) frame.Page {
	i := textinput.New()
	i.Placeholder = "Enter your passphrase"
	i.Focus()
	i.CharLimit = 128
	i.Width = 40
	i.EchoMode = textinput.EchoPassword
	i.EchoCharacter = '•'

	choices := []frame.Choice{
		{Label: "Unlock", Value: "unlock"},
		{Label: "Switch wallet", Value: "switch"},
		{Label: "Create new wallet", Value: "create"},
	}

	return &state{
		ctx:      ctx,
		input:    i,
		selector: frame.NewChoiceSelector(choices),
	}
}

func (m *state) Init() tea.Cmd {
	// Ensure text input is focused
	m.input.Focus()
	return textinput.Blink
}

func (m *state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		// Handle quit keys first
		if v.Type == tea.KeyCtrlC || v.Type == tea.KeyEsc {
			return m, tea.Quit
		}

		// Handle special shortcuts
		if v.Type == tea.KeyCtrlS && v.String() == "ctrl+s" {
			return m, frame.Navigate(config.SwitchWalletPage)
		}

		// Check if this is a navigation key that should be handled by choice selector
		if v.String() == "up" || v.String() == "down" || v.String() == "k" || v.String() == "j" || v.String() == "enter" {
			choiceResult := m.selector.Update(msg)
			if choiceResult.Consumed {
				// Handle choice selection
				if choiceResult.Action == frame.ActionSelection && choiceResult.Selected != nil {
					switch choiceResult.Selected.Value {
					case "unlock":
						pass := m.input.Value()
						if err := m.ctx.Router.Unlock(pass); err != nil {
							m.input.SetValue("")
							m.err = err.Error()
							return m, nil
						}
						m.ctx.TempPassphrase = pass
						return m, frame.Navigate(config.HomePage)
					case "switch":
						m.ctx.TempPassphrase = ""
						m.input.SetValue("")
						return m, frame.Navigate(config.SwitchWalletPage)
					case "create":
						return m, frame.Navigate(config.CreateWalletPage)
					}
				}
				return m, nil
			}
		}
	}

	// Let text input handle all other messages
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *state) View() string {
	v := frame.NewViewBuilder()
	name, _ := m.ctx.WalletRepo.GetActiveWalletName()
	v.Line(fmt.Sprintf("Enter passphrase to unlock wallet %s\n", color.New(color.Bold).Sprintf("%s", name)))
	v.Line(m.input.View())
	v.WithErrText(m.err)
	v.Line("")
	v.Line(m.selector.RenderWithPrefix("▶"))
	return v.Build()
}
