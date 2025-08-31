package sync

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/satelliondao/satellion/ui/frame"
	"github.com/satelliondao/satellion/wallet"
)

type BalanceState int

const (
	BalanceIdle BalanceState = iota
	BalanceScanning
	BalanceComplete
	BalanceError
)

type balanceState struct {
	ctx        *frame.AppContext
	state      BalanceState
	info       *wallet.BalanceInfo
	err        error
	progress   float64
	onComplete func(*wallet.BalanceInfo, error)
}

type balanceCompleteMsg struct {
	info *wallet.BalanceInfo
	err  error
}

type balanceProgressMsg struct {
	progress float64
}

func NewBalanceComponent(ctx *frame.AppContext) *balanceState {
	return &balanceState{
		ctx:   ctx,
		state: BalanceIdle,
	}
}

func (bc *balanceState) SetOnComplete(callback func(*wallet.BalanceInfo, error)) {
	bc.onComplete = callback
}

func (bc *balanceState) StartScan() tea.Cmd {
	if bc.state == BalanceScanning {
		return nil
	}
	bc.state = BalanceScanning
	bc.err = nil
	bc.info = nil
	bc.progress = 0
	return bc.scanBalance()
}

func (bc *balanceState) Update(msg tea.Msg) tea.Cmd {
	switch v := msg.(type) {
	case balanceCompleteMsg:
		bc.state = BalanceComplete
		bc.info = v.info
		bc.err = v.err
		bc.progress = 100
		if v.err != nil {
			bc.state = BalanceError
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

func (bc *balanceState) IsScanning() bool {
	return bc.state == BalanceScanning
}

func (bc *balanceState) IsComplete() bool {
	return bc.state == BalanceComplete
}

func (bc *balanceState) HasError() bool {
	return bc.state == BalanceError
}

func (bc *balanceState) GetInfo() *wallet.BalanceInfo {
	return bc.info
}

func (bc *balanceState) GetError() error {
	return bc.err
}

func (bc *balanceState) View() string {
	v := frame.NewViewBuilder().HideLogo()
	switch bc.state {
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

func (bc *balanceState) scanBalance() tea.Cmd {
	return func() tea.Msg {
		info, err := bc.ctx.Router.GetWalletBalanceInfo(bc.ctx.TempPassphrase)
		return balanceCompleteMsg{info: info, err: err}
	}
}

func (bc *balanceState) SendProgress(progress float64) tea.Cmd {
	return func() tea.Msg {
		return balanceProgressMsg{progress: progress}
	}
}
