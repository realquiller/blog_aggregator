package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(homePath, configFileName)
	return filePath, nil
}

func Read() (Config, error) {
	// Read the file
	filePath, err := getConfigFilePath()

	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(filePath)

	if err != nil {
		return Config{}, err
	}

	// Parse JSON into struct
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func write(cfg Config) error {
	// get the file path
	filePath, err := getConfigFilePath()

	if err != nil {
		return err
	}

	// convert struct to JSON

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(*cfg)
}
