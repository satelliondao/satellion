package usecase

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/cfg"
	"github.com/satelliondao/satellion/chain"
)

type head struct {
	Height    int
	Timestamp time.Time
}

type syncModel struct {
	ch       *chain.Chain
	peers    int
	minPeers int
	inSync   bool
	head     head
}

type tickMsg time.Time
type syncCompleteMsg struct {
	head head
}

func (m syncModel) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m syncModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case syncCompleteMsg:
		m.inSync = true
		m.head = msg.head
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tickMsg:
		stamp, err := m.ch.BestBlock()
		if err != nil {
			return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
		}
		m.peers = int(m.ch.ConnectedCount())
		m.head = head{Height: int(stamp.Height), Timestamp: stamp.Timestamp}
		isCurrent := time.Since(stamp.Timestamp) < 10*time.Minute

		if isCurrent && m.peers >= m.minPeers {
			m.inSync = false
			return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return syncCompleteMsg{head: m.head}
			})
		}
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
	}
	return m, nil
}

func (m syncModel) View() string {
	if m.inSync {
		return color.New(color.FgYellow).Sprintf("Syncing blockchain...\nBest height=%d \ntime=%s \npeers=%d\nPress q to quit.\n", m.head.Height, m.head.Timestamp.UTC().Format(time.RFC3339), m.peers)
	}

	return color.New(color.FgGreen).Sprintf("Synchronization complete, head at height %d\nPress q to quit.\n", m.head.Height)
}

func (wm *Router) Sync() {
	loaded, err := cfg.Load()
	if err != nil {
		fmt.Println("failed to load config:", err)
		os.Exit(1)
	}
	ch := chain.NewChain(loaded)
	if err := ch.Start(); err != nil {
		fmt.Println("failed to start chain service:", err)
		os.Exit(1)
	}
	defer func() { _ = ch.Stop() }()
	m := syncModel{ch: ch, minPeers: loaded.MinPeers}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("sync UI error:", err)
		os.Exit(1)
	}
}
