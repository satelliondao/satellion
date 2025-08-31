package send

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/page"
)

type state struct {
	ctx *framework.AppContext
}

func New(ctx *framework.AppContext) framework.Page {
	return &state{ctx: ctx}
}

func (m *state) Init() tea.Cmd {
	return nil
}

func (m *state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		if stdout.ShouldQuit(v) {
			return m, tea.Quit
		}
		if v.Type == tea.KeyEsc {
			return m, framework.Navigate(page.Home)
		}
	}
	return m, nil
}

func (m *state) View() string {
	v := framework.NewViewBuilder()
	v.Line("Send")
	return v.Build()
}
