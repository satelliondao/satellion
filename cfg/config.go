package cfg

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
  Peers  []string `json:"peers"`
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
	return &config, nil	
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


