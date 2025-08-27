package usecase

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	prompt "github.com/satelliondao/satellion/cli/promt"
	"github.com/satelliondao/satellion/stdout"
)

func (wm *Router) ImportWalletFromSeed() {
	reader := bufio.NewReader(os.Stdin)
	mnemonic, err := prompt.ProvideMnemonic()
	if err != nil {
		stdout.Error.Printf("Failed to read mnemonic: %v\n", err)
		return
	}

	wallet, err := createHDWalletFromSeed(mnemonic)
	if err != nil {
		stdout.Error.Printf("Failed to create HD wallet from seed phrase: %v\n", err)
		return
	}

	fmt.Print("Enter a name for this wallet: ")
	walletName, _ := reader.ReadString('\n')
	walletName = strings.TrimSpace(walletName)

	if walletName == "" {
		walletName = "Imported HD Wallet " + time.Now().Format("2006-01-02 15:04:05")
	}

	err = wm.WalletRepo.Add(walletName, wallet.Mnemonic)
	if err != nil {
		stdout.Error.Printf("Failed to add wallet to list: %v\n", err)
		return
	}

	DisplayWalletInfo(wallet)
	stdout.Success.Printf("HD wallet '%s' imported and saved securely!", walletName)
}
