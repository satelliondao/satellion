package usecase

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/satelliondao/satellion/cli/palette"
	prompt "github.com/satelliondao/satellion/cli/promt"
	"github.com/satelliondao/satellion/ports"
)

func (wm *Router) ImportWalletFromSeed() {
	reader := bufio.NewReader(os.Stdin)
	mnemonic, err := prompt.ProvideMnemonic()
	if err != nil {
		palette.Error.Printf("Failed to read mnemonic: %v\n", err)
		return
	}

	hdWallet, err := createHDWalletFromSeed(string(mnemonic))
	if err != nil {
		palette.Error.Printf("Failed to create HD wallet from seed phrase: %v\n", err)
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
		palette.Error.Printf("Failed to add wallet to list: %v\n", err)
		return
	}

	DisplayWalletInfo(hdWallet)
	palette.Success.Printf("HD wallet '%s' imported and saved securely!", walletName)
}
