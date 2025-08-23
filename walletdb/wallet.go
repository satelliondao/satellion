package walletdb

type WalletEntity struct {
	Name     string `json:"name"`
	Mnemonic string `json:"mnemonic"`
}

type WalletStore struct {	
	Wallets []WalletEntity
}