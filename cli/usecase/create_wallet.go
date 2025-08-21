package usecase

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/satelliondao/satellion/cli/palette"
	prompt "github.com/satelliondao/satellion/cli/promt"
	"github.com/satelliondao/satellion/utils/term"
)

func (wm *Router) GenerateNewWallet() {
	palette.Info.Println("Generating new master key")
	wallet := genNewWallet()
	palette.Info.Print("Enter wallet name: ")
	reader := bufio.NewReader(os.Stdin)
	walletName, _ := reader.ReadString('\n')
	walletName = strings.TrimSpace(walletName)
	if walletName == "" {
		palette.Error.Println("Wallet name is required")
		return
	}
	wallet.Name = walletName
	palette.Error.Printf("🔑 %s", wallet.SeedPhrase)
	term.Newline()
	fmt.Println("Make sure to write down your seed phrase in a safe place")
	palette.Warning.Println("Press enter to continue")

	// wait while user press enter
	reader = bufio.NewReader(os.Stdin)
	reader.ReadString('\n')


	if !prompt.VerifyMnemonicSaved(wallet.SeedPhrase) {
		palette.Error.Println("Mnemonic verification failed. Aborting.")
		return
	}
}
