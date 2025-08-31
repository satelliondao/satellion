package neutrino

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/gcs/builder"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/satelliondao/satellion/config"
	"github.com/stretchr/testify/assert"
)

const TestBlockHeight = 879614

const TestAddress = "bc1pj587y3psgrlsyfsqzmgsy6yun2atpgkwzu03e4lfhm6a2juqchdqyd2g45"

func TestGetCompactFilter(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	cnains := NewChainService(cfg)
	err = cnains.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer cnains.Stop()
	blockhash, err := cnains.GetBlockHash(TestBlockHeight)
	if err != nil {
		t.Fatal(err)
	}
	filter, err := cnains.neutrino.GetCFilter(*blockhash, wire.GCSFilterRegular)
	if err != nil {
		t.Fatal(err)
	}
	address, err := btcutil.DecodeAddress("bc1pj587y3psgrlsyfsqzmgsy6yun2atpgkwzu03e4lfhm6a2juqchdqyd2g45", &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	script, err := txscript.PayToAddrScript(address)
	if err != nil {
		t.Fatal(err)
	}
	key := builder.DeriveKey(blockhash)
	match, err := filter.Match(key, script)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, match, "Should match")
}
