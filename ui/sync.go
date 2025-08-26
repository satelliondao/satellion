package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/cfg"
	"github.com/satelliondao/satellion/stdout"
)

type head struct {
	Height    int
	Timestamp time.Time
}

type tickMsg time.Time
type completeMsg struct {
	head head
}
type errorMsg struct {
	err error
}

type syncModel struct {
	ctx        *AppContext
	head       head
	peers      int
	isComplete bool
	err        error
}

func NewSync(ctx *AppContext) Page { return &syncModel{ctx: ctx} }

func (s *syncModel) Init() tea.Cmd {
	if err := s.ctx.Router.StartChain(); err != nil {
		return func() tea.Msg { return errorMsg{err: err} }
	}
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (s *syncModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		if stdout.ShouldQuit(v) {
			return s, Navigate(cfg.HomePage)
		}
	case tickMsg:
		stamp, peers, err := s.ctx.Router.BestBlock()
		if err != nil {
			return s, tea.Tick(time.Second, func(t time.Time) tea.Msg { return errorMsg{err: err} })
		}
		s.peers = peers
		isCurrent := time.Since(stamp.Timestamp) < 20*time.Minute
		s.head = head{Height: int(stamp.Height), Timestamp: stamp.Timestamp}
		if isCurrent && peers >= s.ctx.Router.MinPeers() {
			_ = s.ctx.Router.StopChain()
			s.isComplete = true
			return s, nil
		}
		return s, tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
	}
	return s, nil
}

func (s *syncModel) View() string {
	if s.isComplete {
		totalBalance := 0.0
		blockTime := s.head.Timestamp.UTC().Format(time.RFC3339)
		minutesAgo := int(time.Since(s.head.Timestamp).Minutes())
		mempoolURL := ""
		if s.head.Height > 0 {
			mempoolURL = "https://mempool.space/block/" + fmt.Sprintf("%d", s.head.Height)
		}
		return color.New(color.FgGreen).Sprintf(
			"Synchronization complete!\n"+
				"Head at height: %d\n"+
				"Block time:     %s (%d min ago)\n"+
				"Total balance:  %.8f BTC\n"+
				"Explore block:  %s\n\n"+
				stdout.Quit(),
			s.head.Height,
			blockTime,
			minutesAgo,
			totalBalance,
			mempoolURL,
		)
	}

	return color.New(color.FgYellow).Sprintf(
		"⏳ Syncing blockchain...\n"+
			"Best height: %d\n"+
			"Timestamp:   %s\n"+
			"Peers:       %d\n\n"+
			stdout.Quit(),
		s.head.Height,
		s.head.Timestamp.UTC().Format(time.RFC3339),
		s.peers,
	)
}
