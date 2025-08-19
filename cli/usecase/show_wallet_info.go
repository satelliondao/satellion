package usecase

import (
	"fmt"
)

func (wm *Router) ShowWalletInfo() {
	fmt.Print(`
👁️  HD Wallet Information
=========================
`)
	walletList, err := wm.walletRepo.LoadWalletList()
	if err != nil {
		fmt.Printf("❌ Failed to load wallet list: %v\n", err)
		return
	}

	if len(walletList.Wallets) == 0 {
		fmt.Println("❌ No wallets found!")
		fmt.Println("Run 'satellion init' to create or import a wallet.")
		return
	}

	fmt.Println("📋 Available HD Wallets:")
	for i, wallet := range walletList.Wallets {
		defaultIndicator := ""
		if wallet.IsDefault {
			defaultIndicator = " (Default)"
		}
		fmt.Printf(`
%d. %s%s
	Master Address: %s
	Next Address Index: %d
	Used Addresses: %d
	Created: %s
	`, i+1, wallet.Name, defaultIndicator, wallet.Address, wallet.NextIndex, len(wallet.UsedIndexes), wallet.CreatedAt)
	}

	if walletList.Default != "" {
		fmt.Println("🔑 Default HD Wallet Details:")

		hdWallet, err := wm.walletRepo.LoadHDWallet(walletList.Default)
		if err != nil {
			fmt.Printf("❌ Failed to load default wallet: %v\n", err)
			return
		}

		fmt.Printf(`
Master Address: %s
Next Address Index: %d
Used Addresses: %d
`, hdWallet.MasterAddress, hdWallet.NextIndex, len(hdWallet.UsedIndexes))

		nextAddress, err := wm.walletRepo.GetNextAddress(walletList.Default)
		if err == nil {
			fmt.Printf("Next Address: %s\n", nextAddress)
		}
	}
}
