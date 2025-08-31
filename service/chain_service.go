package service

import (
	"fmt"

	"github.com/lightninglabs/neutrino/headerfs"
	"github.com/satelliondao/satellion/config"
	"github.com/satelliondao/satellion/neutrino"
	"github.com/satelliondao/satellion/wallet"
)

type ChainService struct {
	neutrino *neutrino.ChainService
	config   *config.Config
}

func NewChainService(config *config.Config) *ChainService {
	return &ChainService{config: config}
}

func (s *ChainService) Start() error {
	if s.neutrino != nil {
		return nil
	}
	if s.config == nil {
		loaded, err := config.Load()
		if err != nil {
			return err
		}
		s.config = loaded
	}
	s.neutrino = neutrino.NewChainService(s.config)
	return s.neutrino.Start()
}

func (s *ChainService) Stop() error {
	if s.neutrino == nil {
		return nil
	}
	err := s.neutrino.Stop()
	s.neutrino = nil
	return err
}

func (s *ChainService) BestBlock() (*headerfs.BlockStamp, int, error) {
	if s.neutrino == nil {
		return nil, 0, fmt.Errorf("chain not started")
	}
	stamp, err := s.neutrino.BestBlock()
	if err != nil {
		return nil, 0, err
	}
	return stamp, int(s.neutrino.ConnectedCount()), nil
}

func (s *ChainService) GetBalance(w *wallet.Wallet) (*neutrino.BalanceInfo, error) {
	if w == nil {
		return nil, fmt.Errorf("no active wallet")
	}
	balanceService := neutrino.NewBalanceService(s.neutrino)
	return balanceService.ScanLedger(w)
}
