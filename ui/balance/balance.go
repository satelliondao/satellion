package balance

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/neutrino"
	"github.com/satelliondao/satellion/ui/staff"
)

type Status int

const (
	BalanceIdle Status = iota
	BalanceScanning
	BalanceComplete
	BalanceError
)

type State struct {
	ctx        *staff.AppContext
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

func New(ctx *staff.AppContext) *State {
	return &State{
		ctx:    ctx,
		status: BalanceIdle,
	}
}

func (bc *State) SetOnComplete(callback func(*neutrino.BalanceInfo, error)) {
	bc.onComplete = callback
}

func (bc *State) StartScan() tea.Cmd {
	if bc.status == BalanceScanning {
		return nil
	}
	bc.status = BalanceScanning
	bc.err = nil
	bc.info = nil
	bc.progress = 0
	return bc.scanBalance()
}

func (bc *State) Update(msg tea.Msg) tea.Cmd {
	switch v := msg.(type) {
	case balanceCompleteMsg:
		bc.status = BalanceComplete
		bc.info = v.info
		bc.err = v.err
		bc.progress = 100
		if v.err != nil {
			bc.status = BalanceError
		}
		if bc.onComplete != nil {
			bc.onComplete(v.info, v.err)
		}
		return nil
	case balanceProgressMsg:
		bc.progress = v.progress
		return nil
	}
	return nil
}

func (bc *State) IsScanning() bool {
	return bc.status == BalanceScanning
}

func (bc *State) IsComplete() bool {
	return bc.status == BalanceComplete
}

func (bc *State) HasError() bool {
	return bc.status == BalanceError
}

func (bc *State) GetInfo() *neutrino.BalanceInfo {
	return bc.info
}

func (bc *State) GetError() error {
	return bc.err
}

func (bc *State) View() string {
	v := staff.NewViewBuilder().HideLogo()
	switch bc.status {
	case BalanceScanning:
		v.Line(fmt.Sprintf("Scanning... %.1f%%", bc.progress))
	case BalanceError:
		v.Line(color.New(color.FgRed).Sprintf("Error: %v", bc.err))
	case BalanceComplete:
		if bc.info != nil {
			v.Line(fmt.Sprintf("%d sats, %d UTXOs", bc.info.Balance, bc.info.UtxoCount))
		} else {
			v.Line(color.New(color.FgRed).Sprintf("No balance information available"))
		}
	default:
		v.Line(color.New(color.FgHiBlack).Sprintf("Balance not loaded"))
	}
	return v.Build()
}

func (bc *State) scanBalance() tea.Cmd {
	return func() tea.Msg {
		info, err := bc.ctx.Router.GetWalletBalanceInfo(bc.ctx.TempPassphrase)
		return balanceCompleteMsg{info: info, err: err}
	}
}

func (bc *State) SendProgress(progress float64) tea.Cmd {
	return func() tea.Msg {
		return balanceProgressMsg{progress: progress}
	}
}
