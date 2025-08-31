package send

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/page"
	"github.com/satelliondao/satellion/ui/staff"
)

type state struct {
	ctx *staff.AppContext
}

func New(ctx *staff.AppContext) staff.Page {
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
			return m, staff.Navigate(page.Home)
		}
	}
	return m, nil
}

func (m *state) View() string {
	v := staff.NewViewBuilder()
	v.Line("Send")
	return v.Build()
}
