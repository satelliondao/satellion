package persistence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/satelliondao/satellion/enclave"
	"github.com/satelliondao/satellion/ports"
)


type HDWalletRepo struct {
	enclave *enclave.Enclave
}

func NewHDWalletRepo() *HDWalletRepo {
	return &HDWalletRepo{
		enclave: enclave.NewEnclave("hd-wallet-keys"),
	}
}


func (wm *HDWalletRepo) LoadHDWallet(walletID string) (*ports.HDWallet, error) {
	hdWalletData, err := wm.enclave.Load(walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to load HD wallet: %w", err)
	}

	var hdWallet ports.HDWallet
	err = json.Unmarshal(hdWalletData, &hdWallet)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HD wallet: %w", err)
	}

	return &hdWallet, nil
}

func (wm *HDWalletRepo) LoadWalletList() (*ports.WalletList, error) {
	data, err := wm.enclave.Load("wallets.json")
	if err != nil {
		if _, ok := err.(*enclave.NotFoundError); ok {
			return &ports.WalletList{
				Wallets: []ports.WalletInfo{},
				Default: "",
			}, nil
		}
		return nil, fmt.Errorf("failed to load wallet list: %w", err)
	}

	var walletList ports.WalletList
	err = json.Unmarshal(data, &walletList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse wallet list: %w", err)
	}

	return &walletList, nil
}

func (wm *HDWalletRepo) SaveWalletList(walletList *ports.WalletList) error {
	data, err := json.MarshalIndent(walletList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal wallet list: %w", err)
	}

	err = wm.enclave.Save("wallets.json", data)
	if err != nil {
		return fmt.Errorf("failed to save wallet list: %w", err)
	}

	return nil
}

func (wm *HDWalletRepo) AddWallet(walletInfo ports.WalletInfo) error {
	walletList, err := wm.LoadWalletList()
	if err != nil {
		return err
	}

	for _, wallet := range walletList.Wallets {
		if wallet.ID == walletInfo.ID {
			return fmt.Errorf("wallet with ID '%s' already exists", walletInfo.ID)
		}
	}

	if len(walletList.Wallets) == 0 {
		walletInfo.IsDefault = true
		walletList.Default = walletInfo.ID
	}

	walletList.Wallets = append(walletList.Wallets, walletInfo)

	return wm.SaveWalletList(walletList)
}


func (wm *HDWalletRepo) SetDefaultWallet(walletID string) error {
	walletList, err := wm.LoadWalletList()
	if err != nil {
		return err
	}

	walletExists := false
	for i, wallet := range walletList.Wallets {
		if wallet.ID == walletID {
			walletExists = true
			walletList.Wallets[i].IsDefault = true
		} else {
			walletList.Wallets[i].IsDefault = false
		}
	}

	if !walletExists {
		return fmt.Errorf("wallet with ID '%s' not found", walletID)
	}

	walletList.Default = walletID
	return wm.SaveWalletList(walletList)
}

func (wm *HDWalletRepo) GetNextAddress(walletID string) (string, error) {
	hdWalletData, err := wm.enclave.Load(walletID)
	if err != nil {
		return "", fmt.Errorf("failed to load HD wallet: %w", err)
	}

	var hdWallet ports.HDWallet
	err = json.Unmarshal(hdWalletData, &hdWallet)
	if err != nil {
		return "", fmt.Errorf("failed to parse HD wallet: %w", err)
	}

	nextAddress := deriveAddress(hdWallet.MasterPrivateKey, hdWallet.NextIndex)
	
	return nextAddress, nil
}

func (wm *HDWalletRepo) MarkAddressAsUsed(walletID string, addressIndex uint32) error {
	hdWalletData, err := wm.enclave.Load(walletID)
	if err != nil {
		return fmt.Errorf("failed to load HD wallet: %w", err)
	}

	var hdWallet ports.HDWallet
	err = json.Unmarshal(hdWalletData, &hdWallet)
	if err != nil {
		return fmt.Errorf("failed to parse HD wallet: %w", err)
	}

	found := false
	for _, usedIndex := range hdWallet.UsedIndexes {
		if usedIndex == addressIndex {
			found = true
			break
		}
	}
	
	if !found {
		hdWallet.UsedIndexes = append(hdWallet.UsedIndexes, addressIndex)
	}

	if hdWallet.NextIndex == addressIndex {
		hdWallet.NextIndex++
	}

	updatedData, err := json.Marshal(hdWallet)
	if err != nil {
		return fmt.Errorf("failed to marshal HD wallet: %w", err)
	}

	err = wm.enclave.Save(walletID, updatedData)
	if err != nil {
		return fmt.Errorf("failed to save HD wallet: %w", err)
	}

	walletList, err := wm.LoadWalletList()
	if err != nil {
		return err
	}

	for i, wallet := range walletList.Wallets {
		if wallet.ID == walletID {
			walletList.Wallets[i].NextIndex = hdWallet.NextIndex
			walletList.Wallets[i].UsedIndexes = hdWallet.UsedIndexes
			break
		}
	}

	return wm.SaveWalletList(walletList)
}

func (wm *HDWalletRepo) DeleteWallet(walletID string) error {
	walletList, err := wm.LoadWalletList()
	if err != nil {
		return err
	}

	for i, wallet := range walletList.Wallets {
		if wallet.ID == walletID {
			walletList.Wallets = append(walletList.Wallets[:i], walletList.Wallets[i+1:]...)
			break
		}
	}

	return wm.SaveWalletList(walletList)
}

func deriveAddress(masterPrivateKey string, index uint32) string {
	// In real implementation, this would use BIP32 derivation
	// For now, using a simple hash with index
	derivationData := fmt.Sprintf("%s_%d", masterPrivateKey, index)
	hash := sha256.Sum256([]byte(derivationData))
	return "0x" + hex.EncodeToString(hash[:20])
} 