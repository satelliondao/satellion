package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/fatih/color"
	"github.com/satelliondao/satellion/ports"
)

func DisplayHDWalletInfo(hdWallet *ports.HDWallet) {
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

func DeriveAddressFromPrivateKey(privateKey string) string {
	// This would be implemented in the wallet package
	// For now, using a simple hash
	hash := sha256.Sum256([]byte(privateKey))
	return "0x" + hex.EncodeToString(hash[:20])
}
