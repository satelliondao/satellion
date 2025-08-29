package wallet

import (
	"testing"

	"github.com/satelliondao/satellion/mnemonic"
	"github.com/stretchr/testify/assert"
)

func TestDeriveReceiveAddress(t *testing.T) {
	words := []string{
		"begin", "hurt", "asset", "stereo", "remove", "around", "emotion", "interest", "ostrich", "invest", "consider", "extend", "egg", "decorate", "fall",
	}
	testMnemonic := &mnemonic.Mnemonic{Words: words}
	passphrase := "123"
	name := "test-wallet"

	t.Run("BasicDeriveReceiveAddress", func(t *testing.T) {
		wallet := New(testMnemonic, passphrase, name, 1, 1, "")
		initialNextReceiveIndex := wallet.NextReceiveIndex
		address, err := wallet.ReceiveAddress()

		assert.NoError(t, err)
		assert.False(t, address.Change, "Receive address should have Change=false")
		assert.Equal(t, uint32(1), address.DeriviationIndex)
		assert.Equal(t, initialNextReceiveIndex+1, wallet.NextReceiveIndex)
		assert.Equal(t, "xprv9s21ZrQH143K38wjhSDmW83Wy6cAGp9nmKKYwVhF6vqsz1gxFY1RbD5JrwW6JS9wTiwYC4RXfv3Se9B3MEb3XjS9EnBcnndixgTytjEKGEz", wallet.RootKey.String())
		assert.Equal(t, "14LGhiRiye8vS8buAVNydTV7nHq66iy2CM", address.Address)
	})
}
