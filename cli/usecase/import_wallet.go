package usecase

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/satelliondao/satellion/ports"
)

func (wm *Router) ImportWalletFromSeed() {
	fmt.Print(`
📥 Import HD Wallet from Seed Phrase
Enter your 12-word seed phrase:
Seed phrase: `)

	reader := bufio.NewReader(os.Stdin)
	seedPhrase, _ := reader.ReadString('\n')
	seedPhrase = strings.TrimSpace(seedPhrase)

	if seedPhrase == "" {
		fmt.Println("❌ Seed phrase cannot be empty")
		return
	}

	hdWallet, err := createHDWalletFromSeed(seedPhrase)
	if err != nil {
		fmt.Printf("❌ Failed to create HD wallet from seed phrase: %v\n", err)
		return
	}

	fmt.Print("Enter a name for this wallet: ")
	walletName, _ := reader.ReadString('\n')
	walletName = strings.TrimSpace(walletName)

	if walletName == "" {
		walletName = "Imported HD Wallet " + time.Now().Format("2006-01-02 15:04:05")
	}

	walletInfo := ports.WalletInfo{
		ID:          hdWallet.MasterAddress, // Use master address as ID
		Name:        walletName,
		Address:     hdWallet.MasterAddress,
		CreatedAt:   time.Now().Format(time.RFC3339),
		IsDefault:   false,
		NextIndex:   hdWallet.NextIndex,
		UsedIndexes: hdWallet.UsedIndexes,
	}

	err = wm.walletRepo.AddWallet(walletInfo)
	if err != nil {
		fmt.Printf("❌ Failed to add wallet to list: %v\n", err)
		return
	}

	DisplayHDWalletInfo(hdWallet)
	fmt.Printf(`
✅ HD wallet '%s' imported and saved securely!
🔄 Each transaction will use a new address for enhanced privacy.
`, walletName)
}
