package router

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/mnemonic"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/page"
)

type VerifyMnemonicProps struct {
	WalletName string
	Mnemonic   *mnemonic.Mnemonic
}

func Home() tea.Cmd {
	return framework.Navigate(page.Home)
}

func Sync() tea.Cmd {
	return framework.Navigate(page.Sync)
}

func Receive() tea.Cmd {
	return framework.Navigate(page.Receive)
}

func Send() tea.Cmd {
	return framework.Navigate(page.Send)
}

func VerifyMnemonic(walletName string, mnemonic *mnemonic.Mnemonic) tea.Cmd {
	return framework.NavigateWithParams(page.VerifyMnemonic, &VerifyMnemonicProps{WalletName: walletName, Mnemonic: mnemonic})
}

func Passphrase(walletName string, mnemonic *mnemonic.Mnemonic) tea.Cmd {
	return framework.NavigateWithParams(page.Passphrase, &VerifyMnemonicProps{WalletName: walletName, Mnemonic: mnemonic})
}

func ListWallets() tea.Cmd {
	return framework.Navigate(page.ListWallets)
}

func SwitchWallet() tea.Cmd {
	return framework.Navigate(page.SwitchWallet)
}

func UnlockWallet() tea.Cmd {
	return framework.Navigate(page.UnlockWallet)
}

func CreateWallet() tea.Cmd {
	return framework.Navigate(page.CreateWallet)
}
