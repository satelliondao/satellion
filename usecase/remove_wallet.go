package usecase

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func (wm *Router) RemoveWallet() {
	red := color.New(color.FgRed)
	red.Println("🗑️ Remove HD Wallet")
	walletList, err := wm.walletRepo.GetAll()
	if err != nil {
		fmt.Printf("❌ Failed to load wallet list: %v\n", err)
		return
	}
	if len(walletList) == 0 {
		fmt.Println("❌ No wallets to remove.")
		return
	}
	red.Println("Available wallets:")
	for i, wallet := range walletList {
		defaultIndicator := ""
		if wallet.IsDefault {
			defaultIndicator = " (Default)"
		}
		red.Printf("%d. %s%s\n", i+1, wallet.Name, defaultIndicator)
	}

	red.Print("Enter the number of the wallet to remove: ")
	reader := bufio.NewReader(os.Stdin)
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	var choice int
	_, err = fmt.Sscanf(choiceStr, "%d", &choice)
	if err != nil || choice < 1 || choice > len(walletList) {
		red.Println("❌ Invalid choice.")
		return
	}
	selectedWallet := walletList[choice-1]
	red.Printf("Are you sure you want to remove HD wallet '%s'? (y/N): ", selectedWallet.Name)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.ToLower(strings.TrimSpace(confirm))
	if confirm != "y" && confirm != "yes" {
		red.Println("❌ Operation cancelled.")
		return
	}
	err = wm.walletRepo.Delete(selectedWallet.Name)
	if err != nil {
		red.Printf("❌ Failed to remove wallet: %v\n", err)
		return
	}
	red.Printf("✅ HD wallet '%s' removed successfully.\n", selectedWallet.Name)
}
