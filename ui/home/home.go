package home

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/frame"
	"github.com/satelliondao/satellion/ui/frame/page"
	"github.com/satelliondao/satellion/wallet"
)

type state struct {
	ctx      *frame.AppContext
	selector *frame.ChoiceSelector
	w        *wallet.Wallet
	items    []menuItem
}

type menuItem struct{ label, page string }
type errorMsg struct {
	err error
}

var baseMenuItems = []menuItem{
	{label: "Receive", page: page.Receive},
	{label: "Send", page: page.Send},
	{label: "Sync chain", page: page.Sync},
}

func New(ctx *frame.AppContext) frame.Page {
	m := &state{ctx: ctx}
	m.rebuildMenu()
	return m
}

func (m *state) rebuildMenu() {
	items := make([]menuItem, 0, len(baseMenuItems)+1)
	items = append(items, baseMenuItems...)
	m.items = items

	choices := make([]frame.Choice, len(items))
	for i, item := range items {
		choices[i] = frame.Choice{Label: item.label, Value: item.page}
	}

	if m.selector == nil {
		m.selector = frame.NewChoiceSelector(choices)
	} else {
		m.selector.SetChoices(choices)
	}
}

func (m *state) Init() tea.Cmd {
	wallet, err := m.ctx.Router.WalletRepo.GetActiveWallet(m.ctx.TempPassphrase)
	if err != nil {
		return func() tea.Msg { return errorMsg{err: err} }
	}
	m.w = wallet
	return nil
}

func (m *state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	result := m.selector.Update(msg)
	// Handle selection results
	if result.Action == frame.ActionSelection && result.Selected != nil {
		return m, frame.Navigate(result.Selected.Value.(string))
	}
	// Handle other key messages if not consumed by selector
	if !result.Consumed {
		switch v := msg.(type) {
		case tea.KeyMsg:
			if stdout.ShouldQuit(v) {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m *state) View() string {
	v := frame.NewViewBuilder()
	if m.w != nil {
		v.Line(fmt.Sprintf("Wallet %s\n", color.New(color.Bold).Sprintf("%s", m.w.Name)))
	}
	v.Line(m.selector.Render())
	return v.Build()
}
