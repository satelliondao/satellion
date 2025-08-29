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
	err      string
	selected int
}

func New(ctx *frame.AppContext) frame.Page {
	i := textinput.New()
	i.Focus()
	i.CharLimit = 128
	i.Width = 40
	i.EchoMode = textinput.EchoPassword
	i.EchoCharacter = '•'
	return &state{ctx: ctx, input: i}
}

func (m *state) Init() tea.Cmd { return textinput.Blink }

func (m *state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		switch v.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.selected == 0 {
				pass := m.input.Value()
				if err := m.ctx.Router.Unlock(pass); err != nil {
					m.err = err.Error()
					return m, nil
				}
				m.ctx.TempPassphrase = pass
				return m, frame.Navigate(config.HomePage)
			}
			if m.selected == 1 {
				m.ctx.TempPassphrase = ""
				m.input.SetValue("")
				return m, frame.Navigate(config.SwitchWalletPage)
			}
			if m.selected == 2 {
				return m, frame.Navigate(config.CreateWalletPage)
			}
		case tea.KeyUp:
			if m.selected > 0 {
				m.selected--
			}
			return m, nil
		case tea.KeyDown:
			if m.selected < 2 {
				m.selected++
			}
			return m, nil
		}
		if v.Type == tea.KeyCtrlS && v.String() == "ctrl+s" {
			return m, frame.Navigate(config.SwitchWalletPage)
		}
	}
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
	opt := func(i int, label string) string {
		if m.selected == i {
			return color.New(color.FgHiCyan).Sprintf("▶ %s", label)
		}
		return fmt.Sprintf("  %s", label)
	}
	v.Line("")
	v.Line(opt(0, "Unlock"))
	v.Line(opt(1, "Switch wallet"))
	v.Line(opt(2, "Create new wallet"))
	return v.Build()
}
