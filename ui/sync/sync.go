package sync

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/neutrino"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/balance"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/router"
)

type tickMsg time.Time

type state struct {
	ctx        *framework.AppContext
	height     int32
	timestamp  time.Time
	peers      int
	isComplete bool
	balance    *balance.State
}

func New(ctx *framework.AppContext, params interface{}) framework.Page {
	s := &state{ctx: ctx}
	s.balance = balance.New(ctx)
	s.balance.SetOnComplete(s.onBalanceComplete)
	return s
}

func (s *state) Init() tea.Cmd {
	go (func() {
		if err := s.ctx.ChainService.Syncronize(); err != nil {
			panic(err)
		}
	})()
	return s.tick()
}

func (s *state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch v := msg.(type) {
	case tea.KeyMsg:
		if stdout.ShouldQuit(v) || v.Type == tea.KeyEsc {
			return s, router.Home()
		}
		if s.isComplete && v.String() == "r" {
			return s, s.balance.StartScan()
		}
	case tickMsg:
		return s, s.handleTick()
	default:
		cmd = s.balance.Update(msg)
	}
	return s, cmd
}

func (s *state) handleTick() tea.Cmd {
	block, err := s.ctx.ChainService.BestBlock()
	if err != nil {
		return s.tick()
	}
	s.peers = block.Peers
	s.height = block.Height
	s.timestamp = block.Timestamp
	if s.isSynced() {
		s.isComplete = true
		return s.balance.StartScan()
	}
	return s.tick()
}

func (s *state) onBalanceComplete(info *neutrino.BalanceInfo, err error) {
	// Balance info is now handled locally by the balance component
}

func (s *state) isSynced() bool {
	syncTimeout := time.Duration(s.ctx.Config.SyncTimeoutMinutes) * time.Minute
	return time.Since(s.timestamp) < syncTimeout && s.peers >= s.ctx.Config.MinPeers
}

func (s *state) tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (s *state) View() string {
	v := framework.View().
		L(color.New(color.FgHiBlue).Sprintf("Blockchain Sync")).
		L("Height: %d", s.height).
		L("Peers: %d", s.peers).
		L("Last block: %s", s.timestamp.Format("15:04:05")).
		L("")
	if s.isComplete {
		v.L(color.New(color.FgGreen).Sprintf("✓ Synced")).
			L(s.balance.View()).
			Help("R to rescan")
	} else {
		v.Warn("⏳ Syncing...")
	}
	return v.QuitHint().Build()
}
