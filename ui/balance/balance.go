package balance

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/neutrino"
	"github.com/satelliondao/satellion/ui/framework"
)

type Status int

const (
	BalanceIdle Status = iota
	BalanceScanning
	BalanceComplete
	BalanceError
)

type State struct {
	ctx        *framework.AppContext
	status     Status
	info       *neutrino.BalanceInfo
	err        error
	progress   float64
	onComplete func(*neutrino.BalanceInfo, error)
}

type balanceCompleteMsg struct {
	info *neutrino.BalanceInfo
	err  error
}

type balanceProgressMsg struct {
	progress float64
}

func New(ctx *framework.AppContext) *State {
	return &State{
		ctx:    ctx,
		status: BalanceIdle,
	}
}

func (s *State) SetOnComplete(callback func(*neutrino.BalanceInfo, error)) {
	s.onComplete = callback
}

func (s *State) StartScan() tea.Cmd {
	if s.status == BalanceScanning {
		return nil
	}
	s.status = BalanceScanning
	s.err = nil
	s.info = nil
	s.progress = 0
	return s.scanBalance()
}

func (s *State) Update(msg tea.Msg) tea.Cmd {
	switch v := msg.(type) {
	case balanceCompleteMsg:
		s.status = BalanceComplete
		s.info = v.info
		s.err = v.err
		s.progress = 100
		if v.err != nil {
			s.status = BalanceError
		}
		if s.onComplete != nil {
			s.onComplete(v.info, v.err)
		}
		return nil
	case balanceProgressMsg:
		s.progress = v.progress
		return nil
	}
	return nil
}

func (s *State) View() string {
	v := framework.View().HideLogo()
	switch s.status {
	case BalanceScanning:
		v.L("Scanning... %.1f%%", s.progress)
	case BalanceError:
		v.L(color.New(color.FgRed).Sprintf("Error: %v", s.err))
	case BalanceComplete:
		if s.info != nil {
			v.L("%d sats, %d UTXOs", s.info.Balance, s.info.UtxoCount)
		} else {
			v.Err("No balance information available")
		}
	default:
		v.Err("Balance not loaded")
	}
	return v.Build()
}

func (s *State) scanBalance() tea.Cmd {
	return func() tea.Msg {
		wallet, err := s.ctx.WalletRepo.GetActiveWallet(s.ctx.Passphrase)
		if err != nil {
			return balanceCompleteMsg{err: err}
		}

		info, err := neutrino.NewBalance(s.ctx.ChainService).ScanLedger(wallet)
		if err != nil {
			return balanceCompleteMsg{err: err}
		}
		return balanceCompleteMsg{info: info, err: err}
	}
}
