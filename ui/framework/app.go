package framework

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Page interface {
	tea.Model
}

type PageFactory func(ctx *AppContext, i interface{}) Page

type NavigateMsg struct {
	To     string
	Params interface{}
}

func Navigate(to string) tea.Cmd {
	return func() tea.Msg { return NavigateMsg{To: to} }
}

func NavigateWithParams[T any](to string, params T) tea.Cmd {
	return func() tea.Msg { return NavigateMsg{To: to, Params: params} }
}

func HandleNav(msg tea.Msg, home tea.Cmd) tea.Cmd {
	switch v := msg.(type) {
	case tea.KeyMsg:
		if v.Type == tea.KeyCtrlC {
			return tea.Quit
		}
		if v.Type == tea.KeyEsc {
			return home
		}
	}
	return nil
}

type App struct {
	ctx     *AppContext
	pages   map[string]PageFactory
	current Page
}

func NewApp(ctx *AppContext, pages map[string]PageFactory, start string) *App {
	cur := pages[start](ctx, interface{}(nil))
	return &App{ctx: ctx, pages: pages, current: cur}
}

func (a *App) Init() tea.Cmd { return a.current.Init() }

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case NavigateMsg:
		f, ok := a.pages[m.To]
		if !ok {
			return a, nil
		}
		a.current = f(a.ctx, m.Params)
		return a, a.current.Init()
	}
	next, cmd := a.current.Update(msg)
	if p, ok := next.(Page); ok {
		a.current = p
	}
	return a, cmd
}

func (a *App) View() string { return a.current.View() }
