package usecase

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	prompt "github.com/satelliondao/satellion/cli/promt"
	"github.com/satelliondao/satellion/stdout"
	"github.com/satelliondao/satellion/utils/term"
	"github.com/satelliondao/satellion/wallet"
)

type createModel struct {
	wallet *wallet.Wallet
}

func (m createModel) Init() tea.Cmd { return nil }

func (m createModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m createModel) View() string {
	return "Create wallet"
}

func (wm *Router) CreateWallet() {
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
