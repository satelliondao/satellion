package usecase

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/satelliondao/satellion/ports"
)

func (wm *Router) GenerateNewWallet() {
	fmt.Println("🆕 Generating New HD Wallet")
	fmt.Println("===========================")

	// Generate new HD wallet
	hdWallet := generateNewHDWallet()

	// Get wallet name from user
	fmt.Print("Enter a name for this wallet: ")
	reader := bufio.NewReader(os.Stdin)
	walletName, _ := reader.ReadString('\n')
	walletName = strings.TrimSpace(walletName)

	if walletName == "" {
		walletName = "HD Wallet " + time.Now().Format("2006-01-02 15:04:05")
	}

	// Create wallet info
	walletInfo := ports.WalletInfo{
		ID:          hdWallet.MasterAddress, // Use master address as ID
		Name:        walletName,
		Address:     hdWallet.MasterAddress,
		CreatedAt:   time.Now().Format(time.RFC3339),
		IsDefault:   false,
		NextIndex:   hdWallet.NextIndex,
		UsedIndexes: hdWallet.UsedIndexes,
	}

	// Add wallet to list
	err := wm.walletRepo.AddWallet(walletInfo)
	if err != nil {
		fmt.Printf("❌ Failed to add wallet to list: %v\n", err)
		return
	}
	DisplayHDWalletInfo(hdWallet)
	fmt.Printf(`
✅ New HD wallet '%s' generated and saved securely!
💡 Make sure to write down your seed phrase in a safe place.
🔄 Each transaction will use a new address for enhanced privacy.
`, walletName)

}
