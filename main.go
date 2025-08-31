package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/satelliondao/satellion/ui/framework"
	"github.com/satelliondao/satellion/ui/home"
	"github.com/satelliondao/satellion/ui/page"
	"github.com/satelliondao/satellion/ui/passphrase"
	"github.com/satelliondao/satellion/ui/receive"
	"github.com/satelliondao/satellion/ui/send"
	"github.com/satelliondao/satellion/ui/sync"
	"github.com/satelliondao/satellion/ui/verify_mnemonic"
	"github.com/satelliondao/satellion/ui/wallet_create"
	"github.com/satelliondao/satellion/ui/wallet_list"
	"github.com/satelliondao/satellion/ui/wallet_switch"
	"github.com/satelliondao/satellion/ui/wallet_unlock"
)

func main() {
	ctx, err := framework.NewContext()
	if err != nil {
		log.Fatalf("Failed to initialize app context: %v", err)
	}
	pages := map[string]framework.PageFactory{
		page.Home:           home.New,
		page.Sync:           sync.New,
		page.CreateWallet:   wallet_create.New,
		page.VerifyMnemonic: verify_mnemonic.New,
		page.Passphrase:     passphrase.New,
		page.ListWallets:    wallet_list.New,
		page.SwitchWallet:   wallet_switch.New,
		page.UnlockWallet:   wallet_unlock.New,
		page.Receive:        receive.New,
		page.Send:           send.New,
	}
	walletCount, err := ctx.WalletRepo.WalletCount()
	if err != nil {
		panic(err)
	}
	app := framework.NewApp(ctx, pages, startPage(walletCount))
	_, _ = tea.NewProgram(app, tea.WithAltScreen()).Run()
}

func startPage(walletCount int) string {
	if walletCount == 0 {
		return page.CreateWallet
	}
	return page.UnlockWallet
}
