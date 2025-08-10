package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "satellion",
	Short: "Satellion wallet",
	Long: `Satellion wallet is a open source wallet developed by Satellion DAO`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a new wallet with random seed phrase",
	Long: `Generate a new wallet with a cryptographically secure random seed phrase.
The seed phrase will be displayed once - make sure to write it down safely!`,
	Run: func(cmd *cobra.Command, args []string) {
		walletManager := NewWalletManager()
		walletManager.GenerateNewWallet()
	},
}

// ImportCmd imports wallet from seed phrase
var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import wallet from seed phrase",
	Long: `Import an existing wallet using a 12-word seed phrase.
Make sure you're in a secure environment when entering your seed phrase.`,
	Run: func(cmd *cobra.Command, args []string) {
		walletManager := NewWalletManager()
		walletManager.ImportWalletFromSeed()
	},
}

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "List current wallet information",
	Long: `Display the current wallet's address, public key, and seed phrase.
The private key is stored securely and not displayed by default.`,
	Run: func(cmd *cobra.Command, args []string) {
		walletManager := NewWalletManager()
		walletManager.ShowWalletInfo()
	},
}

// ListCmd lists all wallets
var ListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all wallets",
	Long: `Display a list of all available wallets with their names, addresses, and creation dates.
The default wallet is marked with a star.`,
	Run: func(cmd *cobra.Command, args []string) {
		walletManager := NewWalletManager()
		walletManager.ListWallets()
	},
}

// RemoveCmd removes a wallet
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a wallet",
	Long: `Remove a wallet from the system. This will permanently delete the wallet's
private key and remove it from the wallet list. This action cannot be undone.`,
	Run: func(cmd *cobra.Command, args []string) {
		walletManager := NewWalletManager()
		walletManager.RemoveWallet()
	},
}

// DefaultCmd sets default wallet
var DefaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Set default wallet",
	Long: `Set a wallet as the default wallet. The default wallet will be used
for operations when no specific wallet is specified.`,
	Run: func(cmd *cobra.Command, args []string) {
		walletManager := NewWalletManager()
		walletManager.SetDefaultWallet()
	},
}

func SetupCommands() {
	RootCmd.AddCommand(NewCmd)
	RootCmd.AddCommand(ImportCmd)
	RootCmd.AddCommand(ShowCmd)
	RootCmd.AddCommand(ListCmd)
	RootCmd.AddCommand(RemoveCmd)
	RootCmd.AddCommand(DefaultCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 