package usecase

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func (wm *Router) SetDefaultWallet() {
	cyan := color.New(color.FgCyan)
	cyan.Println("⭐ Set Default HD Wallet")
	cyan.Println("=======================")

	walletList, err := wm.WalletRepo.GetAll()
	if err != nil {
		fmt.Printf("❌ Failed to load wallet list: %v\n", err)
		return
	}
	if len(walletList) == 0 {
		fmt.Println("❌ No wallets available.")
		return
	}

	cyan.Println("Available wallets:")
	for i, wallet := range walletList {
		defaultIndicator := ""
		if wallet.IsDefault {
			defaultIndicator = " (Current Default)"
		}
		cyan.Printf("%d. %s%s\n", i+1, wallet.Name, defaultIndicator)
	}
	cyan.Print("Enter the number of the wallet to set as default: ")
	reader := bufio.NewReader(os.Stdin)
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	var choice int
	_, err = fmt.Sscanf(choiceStr, "%d", &choice)
	if err != nil || choice < 1 || choice > len(walletList) {
		cyan.Println("❌ Invalid choice.")
		return
	}
	selectedWallet := walletList[choice-1]

	err = wm.WalletRepo.SetDefault(selectedWallet.Name)
	if err != nil {
		cyan.Printf("❌ Failed to set default wallet: %v\n", err)
		return
	}
	cyan.Printf("✅ HD wallet '%s' set as default.\n", selectedWallet.Name)
}
