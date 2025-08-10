package ports

type WalletInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
	IsDefault bool   `json:"is_default"`
	NextIndex uint32 `json:"next_index"` // Next unused address index
	UsedIndexes []uint32 `json:"used_indexes"` // Indexes of addresses that have been used
}

type WalletList struct {
	Wallets []WalletInfo `json:"wallets"`
	Default string       `json:"default"`
}

type HDWallet struct {
	MasterPrivateKey string   `json:"master_private_key"`
	MasterPublicKey  string   `json:"master_public_key"`
	MasterAddress    string   `json:"master_address"`
	SeedPhrase       string   `json:"seed_phrase"`
	NextIndex        uint32   `json:"next_index"`
	UsedIndexes      []uint32 `json:"used_indexes"`
}


