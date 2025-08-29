package stdout

import (
	tea "github.com/charmbracelet/bubbletea"
)

func ShouldQuit(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyCtrlC
}

func HandleQuit(msg tea.KeyMsg, m tea.Model) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	}
	return m, nil
}
