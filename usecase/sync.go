package usecase

import (
	"fmt"
	"os"

	"github.com/satelliondao/satellion/cfg"
	"github.com/satelliondao/satellion/chain"
)

func (wm *Router) Sync() {
	var loaded *cfg.Config
	loaded, err := cfg.Load()
	if err != nil {
		fmt.Println("failed to load config:", err)
		os.Exit(1)
	}

	ch := chain.NewChain(loaded)
	if err := ch.Sync(); err != nil {
		fmt.Println("failed to start chain service:", err)
		os.Exit(1)
	}
}