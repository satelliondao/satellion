package usecase

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	prompt "github.com/satelliondao/satellion/cli/promt"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/utils/term"
)

func (wm *Router) GenerateNewWallet() {
	stdout.Info.Println("Generating new master key")
	wallet := genNewWallet()
	stdout.Info.Print("Enter wallet name: ")
	reader := bufio.NewReader(os.Stdin)
	walletName, _ := reader.ReadString('\n')
	walletName = strings.TrimSpace(walletName)
	if walletName == "" {
		stdout.Error.Println("Wallet name is required")
		return
	}
	wallet.Name = walletName
	stdout.Error.Printf("🔑 %s", wallet.Mnemonic)
	term.Newline()
	fmt.Println("Make sure to write down your seed phrase in a safe place")
	stdout.Warning.Println("Press enter to continue")

	// wait while user press enter
	reader = bufio.NewReader(os.Stdin)
	reader.ReadString('\n')


	if !prompt.VerifyMnemonicSaved(wallet.Mnemonic) {
		stdout.Error.Println("Mnemonic verification failed. Aborting.")
		return
	}
}
