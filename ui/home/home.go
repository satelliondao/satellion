package home

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/page"
	"github.com/satelliondao/satellion/wallet"
)

type state struct {
	ctx      *framework.AppContext
	selector *framework.ChoiceSelector
	w        *wallet.Wallet
	items    []menuItem
}

type menuItem struct{ label, page string }
type errorMsg struct {
	err error
}

var baseMenuItems = []menuItem{
	{label: "Syncronize blockchain", page: page.Sync},
	{label: "Receive", page: page.Receive},
	{label: "Send", page: page.Send},
}

func New(ctx *framework.AppContext) framework.Page {
	m := &state{ctx: ctx}
	m.rebuildMenu()
	return m
}

func (m *state) rebuildMenu() {
	items := make([]menuItem, 0, len(baseMenuItems)+1)
	items = append(items, baseMenuItems...)
	m.items = items

	choices := make([]framework.Choice, len(items))
	for i, item := range items {
		choices[i] = framework.Choice{Label: item.label, Value: item.page}
	}

	if m.selector == nil {
		m.selector = framework.NewChoiceSelector(choices)
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
	if result.Action == framework.ActionSelection && result.Selected != nil {
		return m, framework.Navigate(result.Selected.Value.(string))
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
	v := framework.NewViewBuilder()
	if m.w != nil {
		v.Line(fmt.Sprintf("Wallet %s\n", color.New(color.Bold).Sprintf("%s", m.w.Name)))
	}
	v.Line(m.selector.Render())
	return v.Build()
}
