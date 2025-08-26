package usecase

import (
	"fmt"
)

func (wm *Router) ListWallets() {
	fmt.Println("📋 HD Wallet List")
	walletList, err := wm.walletRepo.LoadWalletList()
	if err != nil {
		fmt.Printf("❌ Failed to load wallet list: %v\n", err)
		return
	}
	if len(walletList.Wallets) == 0 {
		fmt.Println("No wallets found.")
		return
	}
	for i, wallet := range walletList.Wallets {
		defaultIndicator := ""
		if wallet.IsDefault {
			defaultIndicator = " ⭐"
		}
		fmt.Printf(`
%d. %s%s
	Master Address: %s
	Next Index: %d | Used: %d
	Created: %s
	`, i+1, wallet.Name, defaultIndicator, wallet.Address, wallet.NextIndex, len(wallet.UsedIndexes), wallet.CreatedAt)
	}
}
