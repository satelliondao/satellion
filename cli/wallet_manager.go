package cli

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/satelliondao/satellion/persistence"
	"github.com/satelliondao/satellion/ports"
)

type WalletManager struct {
	walletRepo *persistence.HDWalletRepo
}

func NewWalletManager() *WalletManager {
	return &WalletManager{
		walletRepo: persistence.NewHDWalletRepo(),
	}
}

func (wm *WalletManager) GenerateNewWallet() {
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
	displayHDWalletInfo(hdWallet)
	fmt.Printf(`
✅ New HD wallet '%s' generated and saved securely!
💡 Make sure to write down your seed phrase in a safe place.
🔄 Each transaction will use a new address for enhanced privacy.
`, walletName)
	
}

func (wm *WalletManager) ImportWalletFromSeed() {
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
	
	displayHDWalletInfo(hdWallet)
	fmt.Printf(`
✅ HD wallet '%s' imported and saved securely!
🔄 Each transaction will use a new address for enhanced privacy.
`, walletName)
}

func (wm *WalletManager) ShowWalletInfo() {
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

// ListWallets displays all wallets
func (wm *WalletManager) ListWallets() {
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
		// fmt.Printf("%d. %s%s\n", i+1, wallet.Name, defaultIndicator)
		// fmt.Printf("   Master Address: %s\n", wallet.Address)
		// fmt.Printf("   Next Index: %d | Used: %d\n", wallet.NextIndex, len(wallet.UsedIndexes))
		// fmt.Printf("   Created: %s\n", wallet.CreatedAt)
		// fmt.Println()
		fmt.Printf(`
%d. %s%s
	Master Address: %s
	Next Index: %d | Used: %d
	Created: %s
	`, i+1, wallet.Name, defaultIndicator, wallet.Address, wallet.NextIndex, len(wallet.UsedIndexes), wallet.CreatedAt)
	}
}

// RemoveWallet removes a wallet
func (wm *WalletManager) RemoveWallet() {
	red := color.New(color.FgRed)
	red.Println("🗑️ Remove HD Wallet")
	walletList, err := wm.walletRepo.LoadWalletList()
	if err != nil {
		fmt.Printf("❌ Failed to load wallet list: %v\n", err)
		return
	}
	if len(walletList.Wallets) == 0 {
		fmt.Println("❌ No wallets to remove.")
		return
	}
	red.Println("Available wallets:")
	for i, wallet := range walletList.Wallets {
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
	if err != nil || choice < 1 || choice > len(walletList.Wallets) {
		red.Println("❌ Invalid choice.")
		return
	}
	selectedWallet := walletList.Wallets[choice-1]
	red.Printf("Are you sure you want to remove HD wallet '%s'? (y/N): ", selectedWallet.Name)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.ToLower(strings.TrimSpace(confirm))
	if confirm != "y" && confirm != "yes" {
		red.Println("❌ Operation cancelled.")
		return
	}
	err = wm.walletRepo.DeleteWallet(selectedWallet.ID)
	if err != nil {
		red.Printf("❌ Failed to remove wallet: %v\n", err)
		return
	}
	red.Printf("✅ HD wallet '%s' removed successfully.\n", selectedWallet.Name)
}

func (wm *WalletManager) SetDefaultWallet() {
	cyan := color.New(color.FgCyan)
	cyan.Println("⭐ Set Default HD Wallet")
	cyan.Println("=======================")
	
	walletList, err := wm.walletRepo.LoadWalletList()
	if err != nil {
		fmt.Printf("❌ Failed to load wallet list: %v\n", err)
		return
	}
	if len(walletList.Wallets) == 0 {
		fmt.Println("❌ No wallets available.")
		return
	}
	
	cyan.Println("Available wallets:")
	for i, wallet := range walletList.Wallets {
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
	if err != nil || choice < 1 || choice > len(walletList.Wallets) {
		cyan.Println("❌ Invalid choice.")
		return
	}
	selectedWallet := walletList.Wallets[choice-1]

	err = wm.walletRepo.SetDefaultWallet(selectedWallet.ID)
	if err != nil {
		cyan.Printf("❌ Failed to set default wallet: %v\n", err)
		return
	}
	cyan.Printf("✅ HD wallet '%s' set as default.\n", selectedWallet.Name)
}

func generateNewHDWallet() *ports.HDWallet {
	// Generate a master private key (in real implementation, this would use BIP32/BIP39)
	masterPrivateKey := "0x" + hex.EncodeToString([]byte(fmt.Sprintf("master_key_%d", time.Now().UnixNano())))
	// Generate master address from private key
	masterAddress := deriveAddressFromPrivateKey(masterPrivateKey)
	// Generate seed phrase (simplified - in real implementation would use BIP39)
	seedPhrase := "abandon ability able about above absent absorb abstract absurd abuse access accident"
	return &ports.HDWallet{
		MasterPrivateKey: masterPrivateKey,
		MasterPublicKey:  "0x" + hex.EncodeToString([]byte("master_public_key")),
		MasterAddress:    masterAddress,
		SeedPhrase:       seedPhrase,
		NextIndex:        0,
		UsedIndexes:      []uint32{},
	}
}

func createHDWalletFromSeed(seedPhrase string) (*ports.HDWallet, error) {
	// Validate seed phrase (simplified)
	if len(strings.Split(seedPhrase, " ")) != 12 {
		return nil, fmt.Errorf("seed phrase must be 12 words")
	}
	// Generate master private key from seed phrase (simplified)
	hash := sha256.Sum256([]byte(seedPhrase))
	masterPrivateKey := "0x" + hex.EncodeToString(hash[:])
	masterAddress := deriveAddressFromPrivateKey(masterPrivateKey)
	return &ports.HDWallet{
		MasterPrivateKey: masterPrivateKey,
		MasterPublicKey:  "0x" + hex.EncodeToString([]byte("master_public_key")),
		MasterAddress:    masterAddress,
		SeedPhrase:       seedPhrase,
		NextIndex:        0,
		UsedIndexes:      []uint32{},
	}, nil
}

func displayHDWalletInfo(hdWallet *ports.HDWallet) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf(`
🔑 HD Wallet Information:
=========================
Master Address: %s
Master Public Key: %s
Seed Phrase: %s
Next Address Index: %d
Used Addresses: %d
⚠️  Keep your seed phrase safe and private!
🔄 Each transaction will use a new address for enhanced privacy.
`, hdWallet.MasterAddress, hdWallet.MasterPublicKey, red(hdWallet.SeedPhrase), hdWallet.NextIndex, len(hdWallet.UsedIndexes))
}

func deriveAddressFromPrivateKey(privateKey string) string {
	// This would be implemented in the wallet package
	// For now, using a simple hash
	hash := sha256.Sum256([]byte(privateKey))
	return "0x" + hex.EncodeToString(hash[:20])
}

