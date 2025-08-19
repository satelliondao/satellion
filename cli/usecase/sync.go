package usecase

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/satelliondao/satellion/cfg"
	"github.com/satelliondao/satellion/chain"
)

func (wm *Router) Sync() {
	var cfg cfg.Config
	loaded, err := cfg.Load()
	if err != nil {
		fmt.Println("failed to load config:", err)
		os.Exit(1)
	}

	chainService := chain.NewChainServiceWithPeers(loaded.Peers)
	if err := chainService.Start(); err != nil {
		fmt.Println("failed to start chain service:", err)
		os.Exit(1)
	}
	defer chainService.Stop()

	fmt.Printf("connected peers: %d\n", chainService.Chain.ConnectedCount())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stamp, err := chainService.BestBlock()
			if err != nil {
				fmt.Println("best block error:", err)
				continue
			}
			fmt.Printf("best height=%d time=%s peers=%d\n", stamp.Height, stamp.Timestamp.UTC().Format(time.RFC3339), chainService.Chain.ConnectedCount())
		case <-sigCh:
			fmt.Println("\nshutting down...")
			return
		}
	}
}