package sync

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/ui/frame"
	"github.com/satelliondao/satellion/ui/frame/page"
	"github.com/satelliondao/satellion/wallet"
)

type tickMsg time.Time

type state struct {
	ctx        *frame.AppContext
	height     int
	timestamp  time.Time
	peers      int
	isComplete bool
	balance    *balanceState
}

func New(ctx *frame.AppContext) frame.Page {
	s := &state{ctx: ctx}
	s.balance = NewBalanceComponent(ctx)
	s.balance.SetOnComplete(s.onBalanceComplete)
	return s
}

func (s *state) Init() tea.Cmd {
	if err := s.ctx.Router.StartChain(); err != nil {
		return nil
	}
	return s.tick()
}

func (s *state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch v := msg.(type) {
	case tea.KeyMsg:
		if stdout.ShouldQuit(v) || v.Type == tea.KeyEsc {
			_ = s.ctx.Router.StopChain()
			return s, frame.Navigate(page.Home)
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
	stamp, peers, err := s.ctx.Router.BestBlock()
	if err != nil {
		return s.tick()
	}
	s.peers = peers
	s.height = int(stamp.Height)
	s.timestamp = stamp.Timestamp
	if s.isSynced() {
		s.isComplete = true
		return s.balance.StartScan()
	}
	return s.tick()
}

func (s *state) onBalanceComplete(info *wallet.BalanceInfo, err error) {
	if err == nil && info != nil {
		s.ctx.WalletInfo = info
	}
}

func (s *state) isSynced() bool {
	return time.Since(s.timestamp) < 20*time.Minute && s.peers >= s.ctx.Router.MinPeers()
}

func (s *state) tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (s *state) View() string {
	v := frame.NewViewBuilder()
	if s.isComplete {
		v.Line(s.renderComplete())
		v.Line("")
		v.Line(s.renderBalanceSection())
	} else {
		v.Line(s.renderSyncing())
	}
	return v.Build()
}

func (s *state) renderComplete() string {
	minutesAgo := int(time.Since(s.timestamp).Minutes())
	mempoolURL := fmt.Sprintf("https://mempool.space/block/%d", s.height)
	return fmt.Sprintf(
		color.New(color.FgGreen).Sprintf("Synchronization complete\n")+
			"Block number: %d\n"+
			"Mined %d min ago\n"+
			"Explorer %s",
		s.height, minutesAgo, mempoolURL)
}

func (s *state) renderBalanceSection() string {
	instructions := ""
	if s.balance.IsScanning() {
		instructions = color.New(color.FgHiBlack).Sprintf("Scanning in progress... Press ESC to return to home")
	} else {
		instructions = color.New(color.FgHiBlack).Sprintf("Press 'r' to refresh balance • Press ESC to return to home")
	}
	return fmt.Sprintf(
		color.New(color.FgCyan).Sprintf("Balance\n")+
			"%s\n"+
			"\n"+
			"%s",
		s.balance.View(), instructions)
}

func (s *state) renderSyncing() string {
	return fmt.Sprintf(
		"⏳ Syncing blockchain...\n"+
			"Best height: %d\n"+
			"Timestamp:   %s\n"+
			"Peers:       %d",
		s.height, s.timestamp.UTC().Format(time.RFC3339), s.peers)
}
