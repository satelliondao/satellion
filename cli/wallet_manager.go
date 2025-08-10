package cli

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/satelliondao/satellion/enclave"
)

// WalletInfo represents wallet metadata
type WalletInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"` // This is the master address/ID
	CreatedAt string `json:"created_at"`
	IsDefault bool   `json:"is_default"`
	// HD Wallet specific fields
	NextIndex uint32 `json:"next_index"` // Next unused address index
	UsedIndexes []uint32 `json:"used_indexes"` // Indexes of addresses that have been used
}

// WalletList manages multiple wallets
type WalletList struct {
	Wallets []WalletInfo `json:"wallets"`
	Default string       `json:"default"`
}

// HDWallet represents a hierarchical deterministic wallet
type HDWallet struct {
	MasterPrivateKey string   `json:"master_private_key"`
	MasterPublicKey  string   `json:"master_public_key"`
	MasterAddress    string   `json:"master_address"`
	SeedPhrase       string   `json:"seed_phrase"`
	NextIndex        uint32   `json:"next_index"`
	UsedIndexes      []uint32 `json:"used_indexes"`
}

// WalletManager handles wallet operations
type WalletManager struct {
	enclave *enclave.Enclave
}

// NewWalletManager creates a new wallet manager
func NewWalletManager() *WalletManager {
	return &WalletManager{
		enclave: enclave.NewEnclave("wallet-key"),
	}
}

// loadWalletList loads the wallet list from enclave
func (wm *WalletManager) loadWalletList() (*WalletList, error) {
	data, err := wm.enclave.Load("wallets.json")
	if err != nil {
		if _, ok := err.(*enclave.NotFoundError); ok {
			// Return empty wallet list if file doesn't exist
			return &WalletList{
				Wallets: []WalletInfo{},
				Default: "",
			}, nil
		}
		return nil, fmt.Errorf("failed to load wallet list: %w", err)
	}

	var walletList WalletList
	err = json.Unmarshal(data, &walletList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse wallet list: %w", err)
	}

	return &walletList, nil
}

// saveWalletList saves the wallet list to enclave
func (wm *WalletManager) saveWalletList(walletList *WalletList) error {
	data, err := json.MarshalIndent(walletList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal wallet list: %w", err)
	}

	err = wm.enclave.Save("wallets.json", data)
	if err != nil {
		return fmt.Errorf("failed to save wallet list: %w", err)
	}

	return nil
}

// addWallet adds a new wallet to the list
func (wm *WalletManager) addWallet(walletInfo WalletInfo) error {
	walletList, err := wm.loadWalletList()
	if err != nil {
		return err
	}

	// Check if wallet with same ID already exists
	for _, wallet := range walletList.Wallets {
		if wallet.ID == walletInfo.ID {
			return fmt.Errorf("wallet with ID '%s' already exists", walletInfo.ID)
		}
	}

	// If this is the first wallet, make it default
	if len(walletList.Wallets) == 0 {
		walletInfo.IsDefault = true
		walletList.Default = walletInfo.ID
	}

	walletList.Wallets = append(walletList.Wallets, walletInfo)

	return wm.saveWalletList(walletList)
}

// removeWallet removes a wallet from the list
func (wm *WalletManager) removeWallet(walletID string) error {
	walletList, err := wm.loadWalletList()
	if err != nil {
		return err
	}

	// Find and remove the wallet
	for i, wallet := range walletList.Wallets {
		if wallet.ID == walletID {
			// Remove from slice
			walletList.Wallets = append(walletList.Wallets[:i], walletList.Wallets[i+1:]...)
			
			// If this was the default wallet, set a new default
			if walletList.Default == walletID {
				if len(walletList.Wallets) > 0 {
					walletList.Default = walletList.Wallets[0].ID
					walletList.Wallets[0].IsDefault = true
				} else {
					walletList.Default = ""
				}
			}
			
			// Delete the encrypted HD wallet data
			err = wm.enclave.Delete(walletID)
			if err != nil {
				return fmt.Errorf("failed to delete wallet data: %w", err)
			}
			
			return wm.saveWalletList(walletList)
		}
	}

	return fmt.Errorf("wallet with ID '%s' not found", walletID)
}

// setDefaultWallet sets a wallet as the default
func (wm *WalletManager) setDefaultWallet(walletID string) error {
	walletList, err := wm.loadWalletList()
	if err != nil {
		return err
	}

	// Check if wallet exists
	walletExists := false
	for i, wallet := range walletList.Wallets {
		if wallet.ID == walletID {
			walletExists = true
			// Set this wallet as default
			walletList.Wallets[i].IsDefault = true
		} else {
			// Remove default flag from other wallets
			walletList.Wallets[i].IsDefault = false
		}
	}

	if !walletExists {
		return fmt.Errorf("wallet with ID '%s' not found", walletID)
	}

	walletList.Default = walletID
	return wm.saveWalletList(walletList)
}

// getNextAddress generates the next unused address for a wallet
func (wm *WalletManager) getNextAddress(walletID string) (string, error) {
	// Load HD wallet data
	hdWalletData, err := wm.enclave.Load(walletID)
	if err != nil {
		return "", fmt.Errorf("failed to load HD wallet: %w", err)
	}

	var hdWallet HDWallet
	err = json.Unmarshal(hdWalletData, &hdWallet)
	if err != nil {
		return "", fmt.Errorf("failed to parse HD wallet: %w", err)
	}

	// Generate next address
	nextAddress := deriveAddressFromIndex(hdWallet.MasterPrivateKey, hdWallet.NextIndex)
	
	return nextAddress, nil
}

// markAddressAsUsed marks an address index as used and increments next index
func (wm *WalletManager) markAddressAsUsed(walletID string, addressIndex uint32) error {
	// Load HD wallet data
	hdWalletData, err := wm.enclave.Load(walletID)
	if err != nil {
		return fmt.Errorf("failed to load HD wallet: %w", err)
	}

	var hdWallet HDWallet
	err = json.Unmarshal(hdWalletData, &hdWallet)
	if err != nil {
		return fmt.Errorf("failed to parse HD wallet: %w", err)
	}

	// Add to used indexes if not already there
	found := false
	for _, usedIndex := range hdWallet.UsedIndexes {
		if usedIndex == addressIndex {
			found = true
			break
		}
	}
	
	if !found {
		hdWallet.UsedIndexes = append(hdWallet.UsedIndexes, addressIndex)
	}

	// Update next index if this was the current next index
	if hdWallet.NextIndex == addressIndex {
		hdWallet.NextIndex++
	}

	// Save updated HD wallet
	updatedData, err := json.Marshal(hdWallet)
	if err != nil {
		return fmt.Errorf("failed to marshal HD wallet: %w", err)
	}

	err = wm.enclave.Save(walletID, updatedData)
	if err != nil {
		return fmt.Errorf("failed to save HD wallet: %w", err)
	}

	// Update wallet list metadata
	walletList, err := wm.loadWalletList()
	if err != nil {
		return err
	}

	for i, wallet := range walletList.Wallets {
		if wallet.ID == walletID {
			walletList.Wallets[i].NextIndex = hdWallet.NextIndex
			walletList.Wallets[i].UsedIndexes = hdWallet.UsedIndexes
			break
		}
	}

	return wm.saveWalletList(walletList)
}

// InitializeWallet handles the wallet initialization flow
func (wm *WalletManager) InitializeWallet() {
	fmt.Println("🔐 Wallet Initialization")
	fmt.Println("========================")
	fmt.Println("Choose an option:")
	fmt.Println("1. Generate new HD wallet")
	fmt.Println("2. Import HD wallet from seed phrase")
	fmt.Print("Enter your choice (1 or 2): ")
	
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	
	switch choice {
	case "1":
		wm.GenerateNewWallet()
	case "2":
		wm.ImportWalletFromSeed()
	default:
		fmt.Println("❌ Invalid choice. Please run 'satellion init' again.")
	}
}

// GenerateNewWallet creates a new HD wallet with random seed phrase
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
	walletInfo := WalletInfo{
		ID:          hdWallet.MasterAddress, // Use master address as ID
		Name:        walletName,
		Address:     hdWallet.MasterAddress,
		CreatedAt:   time.Now().Format(time.RFC3339),
		IsDefault:   false,
		NextIndex:   hdWallet.NextIndex,
		UsedIndexes: hdWallet.UsedIndexes,
	}
	
	// Add wallet to list
	err := wm.addWallet(walletInfo)
	if err != nil {
		fmt.Printf("❌ Failed to add wallet to list: %v\n", err)
		return
	}
	
	// Save HD wallet data to encrypted storage
	hdWalletData, err := json.Marshal(hdWallet)
	if err != nil {
		fmt.Printf("❌ Failed to marshal HD wallet: %v\n", err)
		return
	}
	
	err = wm.enclave.Save(hdWallet.MasterAddress, hdWalletData)
	if err != nil {
		fmt.Printf("❌ Failed to save HD wallet: %v\n", err)
		return
	}
	
	// Display wallet information
	displayHDWalletInfo(hdWallet)
	
	fmt.Printf("✅ New HD wallet '%s' generated and saved securely!\n", walletName)
	fmt.Println("💡 Make sure to write down your seed phrase in a safe place.")
	fmt.Println("🔄 Each transaction will use a new address for enhanced privacy.")
}

// ImportWalletFromSeed imports an HD wallet from seed phrase
func (wm *WalletManager) ImportWalletFromSeed() {
	fmt.Println("📥 Import HD Wallet from Seed Phrase")
	fmt.Println("====================================")
	
	// Get seed phrase from user
	fmt.Println("Enter your 12-word seed phrase:")
	fmt.Print("Seed phrase: ")
	
	reader := bufio.NewReader(os.Stdin)
	seedPhrase, _ := reader.ReadString('\n')
	seedPhrase = strings.TrimSpace(seedPhrase)
	
	if seedPhrase == "" {
		fmt.Println("❌ Seed phrase cannot be empty")
		return
	}
	
	// Create HD wallet from seed phrase
	hdWallet, err := createHDWalletFromSeed(seedPhrase)
	if err != nil {
		fmt.Printf("❌ Failed to create HD wallet from seed phrase: %v\n", err)
		return
	}
	
	// Get wallet name from user
	fmt.Print("Enter a name for this wallet: ")
	walletName, _ := reader.ReadString('\n')
	walletName = strings.TrimSpace(walletName)
	
	if walletName == "" {
		walletName = "Imported HD Wallet " + time.Now().Format("2006-01-02 15:04:05")
	}
	
	// Create wallet info
	walletInfo := WalletInfo{
		ID:          hdWallet.MasterAddress, // Use master address as ID
		Name:        walletName,
		Address:     hdWallet.MasterAddress,
		CreatedAt:   time.Now().Format(time.RFC3339),
		IsDefault:   false,
		NextIndex:   hdWallet.NextIndex,
		UsedIndexes: hdWallet.UsedIndexes,
	}
	
	// Add wallet to list
	err = wm.addWallet(walletInfo)
	if err != nil {
		fmt.Printf("❌ Failed to add wallet to list: %v\n", err)
		return
	}
	
	// Save HD wallet data to encrypted storage
	hdWalletData, err := json.Marshal(hdWallet)
	if err != nil {
		fmt.Printf("❌ Failed to marshal HD wallet: %v\n", err)
		return
	}
	
	err = wm.enclave.Save(hdWallet.MasterAddress, hdWalletData)
	if err != nil {
		fmt.Printf("❌ Failed to save HD wallet: %v\n", err)
		return
	}
	
	// Display wallet information
	displayHDWalletInfo(hdWallet)
	
	fmt.Printf("✅ HD wallet '%s' imported and saved securely!\n", walletName)
	fmt.Println("🔄 Each transaction will use a new address for enhanced privacy.")
}

// ShowWalletInfo displays current wallet information
func (wm *WalletManager) ShowWalletInfo() {
	fmt.Println("👁️  HD Wallet Information")
	fmt.Println("=========================")
	
	// Load wallet list
	walletList, err := wm.loadWalletList()
	if err != nil {
		fmt.Printf("❌ Failed to load wallet list: %v\n", err)
		return
	}
	
	if len(walletList.Wallets) == 0 {
		fmt.Println("❌ No wallets found!")
		fmt.Println("Run 'satellion init' to create or import a wallet.")
		return
	}
	
	// Show all wallets
	fmt.Println("📋 Available HD Wallets:")
	fmt.Println("========================")
	
	for i, wallet := range walletList.Wallets {
		defaultIndicator := ""
		if wallet.IsDefault {
			defaultIndicator = " (Default)"
		}
		fmt.Printf("%d. %s%s\n", i+1, wallet.Name, defaultIndicator)
		fmt.Printf("   Master Address: %s\n", wallet.Address)
		fmt.Printf("   Next Address Index: %d\n", wallet.NextIndex)
		fmt.Printf("   Used Addresses: %d\n", len(wallet.UsedIndexes))
		fmt.Printf("   Created: %s\n", wallet.CreatedAt)
		fmt.Println()
	}
	
	// Show default wallet details
	if walletList.Default != "" {
		fmt.Println("🔑 Default HD Wallet Details:")
		fmt.Println("=============================")
		
		// Load HD wallet data
		hdWalletData, err := wm.enclave.Load(walletList.Default)
		if err != nil {
			fmt.Printf("❌ Failed to load default wallet: %v\n", err)
			return
		}
		
		var hdWallet HDWallet
		err = json.Unmarshal(hdWalletData, &hdWallet)
		if err != nil {
			fmt.Printf("❌ Failed to parse HD wallet: %v\n", err)
			return
		}
		
		fmt.Printf("Master Address: %s\n", hdWallet.MasterAddress)
		fmt.Printf("Next Address Index: %d\n", hdWallet.NextIndex)
		fmt.Printf("Used Addresses: %d\n", len(hdWallet.UsedIndexes))
		
		// Show next address
		nextAddress, err := wm.getNextAddress(walletList.Default)
		if err == nil {
			fmt.Printf("Next Address: %s\n", nextAddress)
		}
	}
}

// ListWallets displays all wallets
func (wm *WalletManager) ListWallets() {
	fmt.Println("📋 HD Wallet List")
	fmt.Println("=================")
	
	walletList, err := wm.loadWalletList()
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
		fmt.Printf("%d. %s%s\n", i+1, wallet.Name, defaultIndicator)
		fmt.Printf("   Master Address: %s\n", wallet.Address)
		fmt.Printf("   Next Index: %d | Used: %d\n", wallet.NextIndex, len(wallet.UsedIndexes))
		fmt.Printf("   Created: %s\n", wallet.CreatedAt)
		fmt.Println()
	}
}

// RemoveWallet removes a wallet
func (wm *WalletManager) RemoveWallet() {
	fmt.Println("🗑️  Remove HD Wallet")
	fmt.Println("===================")
	
	walletList, err := wm.loadWalletList()
	if err != nil {
		fmt.Printf("❌ Failed to load wallet list: %v\n", err)
		return
	}
	
	if len(walletList.Wallets) == 0 {
		fmt.Println("❌ No wallets to remove.")
		return
	}
	
	// Show available wallets
	fmt.Println("Available wallets:")
	for i, wallet := range walletList.Wallets {
		defaultIndicator := ""
		if wallet.IsDefault {
			defaultIndicator = " (Default)"
		}
		fmt.Printf("%d. %s%s\n", i+1, wallet.Name, defaultIndicator)
	}
	
	fmt.Print("Enter the number of the wallet to remove: ")
	reader := bufio.NewReader(os.Stdin)
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	
	// Parse choice
	var choice int
	_, err = fmt.Sscanf(choiceStr, "%d", &choice)
	if err != nil || choice < 1 || choice > len(walletList.Wallets) {
		fmt.Println("❌ Invalid choice.")
		return
	}
	
	selectedWallet := walletList.Wallets[choice-1]
	
	// Confirm deletion
	fmt.Printf("Are you sure you want to remove HD wallet '%s'? (y/N): ", selectedWallet.Name)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.ToLower(strings.TrimSpace(confirm))
	
	if confirm != "y" && confirm != "yes" {
		fmt.Println("❌ Operation cancelled.")
		return
	}
	
	// Remove wallet
	err = wm.removeWallet(selectedWallet.ID)
	if err != nil {
		fmt.Printf("❌ Failed to remove wallet: %v\n", err)
		return
	}
	
	fmt.Printf("✅ HD wallet '%s' removed successfully.\n", selectedWallet.Name)
}

// SetDefaultWallet sets a wallet as default
func (wm *WalletManager) SetDefaultWallet() {
	fmt.Println("⭐ Set Default HD Wallet")
	fmt.Println("=======================")
	
	walletList, err := wm.loadWalletList()
	if err != nil {
		fmt.Printf("❌ Failed to load wallet list: %v\n", err)
		return
	}
	
	if len(walletList.Wallets) == 0 {
		fmt.Println("❌ No wallets available.")
		return
	}
	
	// Show available wallets
	fmt.Println("Available wallets:")
	for i, wallet := range walletList.Wallets {
		defaultIndicator := ""
		if wallet.IsDefault {
			defaultIndicator = " (Current Default)"
		}
		fmt.Printf("%d. %s%s\n", i+1, wallet.Name, defaultIndicator)
	}
	
	fmt.Print("Enter the number of the wallet to set as default: ")
	reader := bufio.NewReader(os.Stdin)
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	
	// Parse choice
	var choice int
	_, err = fmt.Sscanf(choiceStr, "%d", &choice)
	if err != nil || choice < 1 || choice > len(walletList.Wallets) {
		fmt.Println("❌ Invalid choice.")
		return
	}
	
	selectedWallet := walletList.Wallets[choice-1]
	
	// Set as default
	err = wm.setDefaultWallet(selectedWallet.ID)
	if err != nil {
		fmt.Printf("❌ Failed to set default wallet: %v\n", err)
		return
	}
	
	fmt.Printf("✅ HD wallet '%s' set as default.\n", selectedWallet.Name)
}

// generateNewHDWallet creates a new HD wallet (simplified implementation)
func generateNewHDWallet() *HDWallet {
	// Generate a master private key (in real implementation, this would use BIP32/BIP39)
	masterPrivateKey := "0x" + hex.EncodeToString([]byte(fmt.Sprintf("master_key_%d", time.Now().UnixNano())))
	
	// Generate master address from private key
	masterAddress := deriveAddressFromPrivateKey(masterPrivateKey)
	
	// Generate seed phrase (simplified - in real implementation would use BIP39)
	seedPhrase := "abandon ability able about above absent absorb abstract absurd abuse access accident"
	
	return &HDWallet{
		MasterPrivateKey: masterPrivateKey,
		MasterPublicKey:  "0x" + hex.EncodeToString([]byte("master_public_key")),
		MasterAddress:    masterAddress,
		SeedPhrase:       seedPhrase,
		NextIndex:        0,
		UsedIndexes:      []uint32{},
	}
}

// createHDWalletFromSeed creates an HD wallet from seed phrase (simplified implementation)
func createHDWalletFromSeed(seedPhrase string) (*HDWallet, error) {
	// Validate seed phrase (simplified)
	if len(strings.Split(seedPhrase, " ")) != 12 {
		return nil, fmt.Errorf("seed phrase must be 12 words")
	}
	
	// Generate master private key from seed phrase (simplified)
	hash := sha256.Sum256([]byte(seedPhrase))
	masterPrivateKey := "0x" + hex.EncodeToString(hash[:])
	
	// Generate master address from private key
	masterAddress := deriveAddressFromPrivateKey(masterPrivateKey)
	
	return &HDWallet{
		MasterPrivateKey: masterPrivateKey,
		MasterPublicKey:  "0x" + hex.EncodeToString([]byte("master_public_key")),
		MasterAddress:    masterAddress,
		SeedPhrase:       seedPhrase,
		NextIndex:        0,
		UsedIndexes:      []uint32{},
	}, nil
}

// displayHDWalletInfo displays HD wallet information
func displayHDWalletInfo(hdWallet *HDWallet) {
	fmt.Println("🔑 HD Wallet Information:")
	fmt.Println("=========================")
	fmt.Printf("Master Address: %s\n", hdWallet.MasterAddress)
	fmt.Printf("Master Public Key: %s\n", hdWallet.MasterPublicKey)
	fmt.Printf("Seed Phrase: %s\n", hdWallet.SeedPhrase)
	fmt.Printf("Next Address Index: %d\n", hdWallet.NextIndex)
	fmt.Printf("Used Addresses: %d\n", len(hdWallet.UsedIndexes))
	fmt.Println("⚠️  Keep your seed phrase safe and private!")
	fmt.Println("🔄 Each transaction will use a new address for enhanced privacy.")
}

// deriveAddressFromPrivateKey derives address from private key
func deriveAddressFromPrivateKey(privateKey string) string {
	// This would be implemented in the wallet package
	// For now, using a simple hash
	hash := sha256.Sum256([]byte(privateKey))
	return "0x" + hex.EncodeToString(hash[:20])
}

// deriveAddressFromIndex derives a new address from master key and index
func deriveAddressFromIndex(masterPrivateKey string, index uint32) string {
	// In real implementation, this would use BIP32 derivation
	// For now, using a simple hash with index
	derivationData := fmt.Sprintf("%s_%d", masterPrivateKey, index)
	hash := sha256.Sum256([]byte(derivationData))
	return "0x" + hex.EncodeToString(hash[:20])
} 