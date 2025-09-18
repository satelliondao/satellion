package send

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/router"
)

type state struct {
	ctx *framework.AppContext
}

func New(ctx *framework.AppContext, params interface{}) framework.Page {
	s := &state{ctx: ctx}
	return s
}

func (m *state) Init() tea.Cmd {
	return nil
}

func (m *state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	nav := framework.HandleNav(msg, router.Home())
	if nav != nil {
		return m, nav
	}

	return m, nil
}

func (m *state) View() string {
	v := framework.View()
	v.L("Send")
	return v.Build()
}
