package usecase

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/satelliondao/satellion/cli/stdout"
	"github.com/satelliondao/satellion/ports"
)

func DisplayWalletInfo(hdWallet *ports.HDWallet) {
stdout.Trace.Printf(`
Next Address Index: %d
Used Addresses: %d
⚠️  Keep your seed phrase safe and private!
🔄 Each transaction will use a new address for enhanced privacy.
`, hdWallet.NextIndex, len(hdWallet.UsedIndexes))
}

func DeriveAddressFromPrivateKey(privateKey string) string {
	// This would be implemented in the wallet package
	// For now, using a simple hash
	hash := sha256.Sum256([]byte(privateKey))
	return "0x" + hex.EncodeToString(hash[:20])
}
