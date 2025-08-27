package stdout

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Quit() string {
	return ("Press q or escape to quit.")
}

func ShouldQuit(msg tea.KeyMsg) bool {
	return msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl+c"
}

func HandleQuit(msg tea.KeyMsg, m tea.Model) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}
