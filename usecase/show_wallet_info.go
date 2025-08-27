package usecase

import (
	"fmt"
)

func (wm *Router) ShowWalletInfo() {
	fmt.Print(`
👁️  HD Wallet Information
=========================
`)
	walletList, err := wm.WalletRepo.GetAll()
	if err != nil {
		fmt.Printf("❌ Failed to load wallet list: %v\n", err)
		return
	}

	if len(walletList) == 0 {
		fmt.Println("❌ No wallets found!")
		fmt.Println("Run 'satellion init' to create or import a wallet.")
		return
	}

	fmt.Println("📋 Available HD Wallets:")
	for i, wallet := range walletList {
		defaultIndicator := ""
		if wallet.IsDefault {
			defaultIndicator = " (Default)"
		}
		fmt.Printf("\n%d. %s%s\n\tNext Address Index: %d\n", i+1, wallet.Name, defaultIndicator, wallet.NextIndex())
	}
}
