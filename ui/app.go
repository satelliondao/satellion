package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Page interface {
	tea.Model
}

type PageFactory func(ctx *AppContext) Page

type NavigateMsg struct{ To string }

func Navigate(to string) tea.Cmd { return func() tea.Msg { return NavigateMsg{To: to} } }

type App struct {
	ctx     *AppContext
	pages   map[string]PageFactory
	current Page
}

func NewApp(ctx *AppContext, pages map[string]PageFactory, start string) *App {
	cur := pages[start](ctx)
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
		a.current = f(a.ctx)
		return a, a.current.Init()
	}
	next, cmd := a.current.Update(msg)
	if p, ok := next.(Page); ok {
		a.current = p
	}
	return a, cmd
}

func (a *App) View() string { return a.current.View() }
