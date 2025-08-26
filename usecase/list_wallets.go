package usecase

import (
	"fmt"
)

func (wm *Router) ListWallets() {
	fmt.Println("📋 HD Wallet List")
	walletList, err := wm.walletRepo.GetAll()
	if err != nil {
		fmt.Printf("❌ Failed to load wallet list: %v\n", err)
		return
	}
	if len(walletList) == 0 {
		fmt.Println("No wallets found.")
		return
	}
	for i, wallet := range walletList {
		defaultIndicator := ""
		if wallet.IsDefault {
			defaultIndicator = " ⭐"
		}
		fmt.Printf("\n%d. %s%s\n\tNext Index: %d\n", i+1, wallet.Name, defaultIndicator, wallet.NextIndex())
	}
}
