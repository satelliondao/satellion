package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	// Peers is an initial list of Bitcoin peers (host:port) that the neutrino client should connect to in the first time.
	Peers []string `json:"peers"`
	// MinPeers is the minimum number of connected peers required before considering sync complete.
	// If omitted or zero in the config file, it defaults to 5.
	MinPeers int `json:"min_peers"`
	// SyncTimeoutMinutes is the maximum age in minutes for a block to be considered current.
	// If omitted or zero in the config file, it defaults to 30 minutes.
	SyncTimeoutMinutes int `json:"sync_timeout_minutes"`
}

func getStoragePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".satellion", "config.json")
}

func (c *Config) Load() (*Config, error) {
	storagePath := getStoragePath()
	jsonFile, err := os.Open(storagePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	var config Config
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	if config.MinPeers == 0 {
		config.MinPeers = 5
	}
	if config.SyncTimeoutMinutes == 0 {
		config.SyncTimeoutMinutes = 30
	}
	return &config, nil
}

// Load reads the configuration from disk using default storage path.
func Load() (*Config, error) {
	var c Config
	return c.Load()
}

func (c *Config) Save(config *Config) error {
	storagePath := getStoragePath()
	if err := os.MkdirAll(filepath.Dir(storagePath), 0o755); err != nil {
		return err
	}
	jsonFile, err := os.Create(storagePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	encoder := json.NewEncoder(jsonFile)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}
	return nil
}
