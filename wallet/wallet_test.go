package wallet

import (
	"fmt"
	"testing"

	"github.com/satelliondao/satellion/mnemonic"
	"github.com/stretchr/testify/assert"
)

func TestBIP86TaprootDerivation(t *testing.T) {
	words := []string{
		"abandon", "abandon", "abandon", "abandon", "abandon", "abandon", "abandon", "abandon", "abandon", "abandon", "abandon", "about",
	}
	testMnemonic := &mnemonic.Mnemonic{Words: words}
	passphrase := ""
	name := "test-wallet"

	t.Run("BIP86TestVectors", func(t *testing.T) {
		wallet := New(testMnemonic, passphrase, name, 0, 0, "")

		assert.Equal(t, "xprv9s21ZrQH143K3GJpoapnV8SFfukcVBSfeCficPSGfubmSFDxo1kuHnLisriDvSnRRuL2Qrg5ggqHKNVpxR86QEC8w35uxmGoggxtQTPvfUu", wallet.RootKey.String())

		firstReceive, err := wallet.ReceiveAddress()
		assert.NoError(t, err)
		assert.False(t, firstReceive.Change, "Receive address should have Change=false")
		assert.Equal(t, uint32(0), firstReceive.DeriviationIndex)
		assert.Equal(t, "bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr", firstReceive.Address.String())

		secondReceive, err := wallet.NewReceiveAddress()
		assert.NoError(t, err)
		assert.False(t, secondReceive.Change, "Receive address should have Change=false")
		assert.Equal(t, uint32(1), secondReceive.DeriviationIndex)
		assert.Equal(t, "bc1p4qhjn9zdvkux4e44uhx8tc55attvtyu358kutcqkudyccelu0was9fqzwh", secondReceive.Address.String())

		firstChange, err := wallet.ChangeAddress()
		assert.NoError(t, err)
		assert.True(t, firstChange.Change, "Change address should have Change=true")
		assert.Equal(t, uint32(0), firstChange.DeriviationIndex)
		assert.Equal(t, "bc1p3qkhfews2uk44qtvauqyr2ttdsw7svhkl9nkm9s9c3x4ax5h60wqwruhk7", firstChange.Address.String())
	})

	t.Run("IndexIncrement", func(t *testing.T) {
		wallet := New(testMnemonic, passphrase, name, 0, 5, "")
		initialReceiveIndex := wallet.NextReceiveIndex
		initialChangeIndex := wallet.NextChangeIndex

		newReceiveAddr, err := wallet.NewReceiveAddress()
		assert.NoError(t, err)
		assert.Equal(t, initialReceiveIndex+1, wallet.NextReceiveIndex, "NextReceiveIndex should be incremented")
		assert.Equal(t, initialReceiveIndex+1, newReceiveAddr.DeriviationIndex)
		assert.False(t, newReceiveAddr.Change)
		assert.True(t, newReceiveAddr.Address.String()[:4] == "bc1p", "Should be a taproot address starting with bc1p")

		newChangeAddr, err := wallet.NewChangeAddress()
		assert.NoError(t, err)
		assert.Equal(t, initialChangeIndex+1, wallet.NextChangeIndex, "NextChangeIndex should be incremented")
		assert.Equal(t, initialChangeIndex+1, newChangeAddr.DeriviationIndex)
		assert.True(t, newChangeAddr.Change)
		assert.True(t, newChangeAddr.Address.String()[:4] == "bc1p", "Should be a taproot address starting with bc1p")
	})

	t.Run("AddressDifferenciation", func(t *testing.T) {
		wallet := New(testMnemonic, passphrase, name, 0, 0, "")

		addr1, err1 := wallet.ReceiveAddress()
		assert.NoError(t, err1)

		addr2, err2 := wallet.NewReceiveAddress()
		assert.NoError(t, err2)

		assert.NotEqual(t, addr1.Address, addr2.Address, "Different indices should produce different addresses")
		assert.Equal(t, uint32(0), addr1.DeriviationIndex)
		assert.Equal(t, uint32(1), addr2.DeriviationIndex)

		changeAddr, err3 := wallet.ChangeAddress()
		assert.NoError(t, err3)
		assert.NotEqual(t, addr1.Address, changeAddr.Address, "Receive and change addresses should be different")
		assert.NotEqual(t, addr2.Address, changeAddr.Address, "Receive and change addresses should be different")
	})

	t.Run("ScriptPubKeyGeneration", func(t *testing.T) {
		wallet := New(testMnemonic, passphrase, name, 0, 0, "")

		address, err := wallet.ReceiveAddress()
		assert.NoError(t, err)
		receiveScriptPubKey, err := address.DeriveTaprootScriptPubKey()
		assert.NoError(t, err)
		assert.Equal(t, 34, len(receiveScriptPubKey), "ScriptPubKey should be 34 bytes")
		assert.Equal(t, byte(0x51), receiveScriptPubKey[0], "First byte should be OP_1 (0x51)")
		assert.Equal(t, byte(0x20), receiveScriptPubKey[1], "Second byte should be 0x20 (32 bytes push)")

		address, err = wallet.ChangeAddress()
		assert.NoError(t, err)
		changeScriptPubKey, err := address.DeriveTaprootScriptPubKey()
		assert.NoError(t, err)
		assert.Equal(t, 34, len(changeScriptPubKey), "ScriptPubKey should be 34 bytes")
		assert.Equal(t, byte(0x51), changeScriptPubKey[0], "First byte should be OP_1 (0x51)")
		assert.Equal(t, byte(0x20), changeScriptPubKey[1], "Second byte should be 0x20 (32 bytes push)")

		assert.NotEqual(t, receiveScriptPubKey, changeScriptPubKey, "Receive and change scriptPubKeys should be different")

		expectedReceiveOutputKey := "a60869f0dbcf1dc659c9cecbaf8050135ea9e8cdc487053f1dc6880949dc684c"
		actualReceiveOutputKey := receiveScriptPubKey[2:]
		assert.Equal(t, expectedReceiveOutputKey, fmt.Sprintf("%x", actualReceiveOutputKey), "Should match BIP 86 test vector output key")
	})
}
