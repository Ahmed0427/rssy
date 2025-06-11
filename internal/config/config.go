package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ConnStr  string `json:"conn_str"`
	Username string `json:"username"`
}

const configFileName = ".rssyconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/" + configFileName, nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (config *Config) Write() error {
	jsonStr, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, []byte(jsonStr), 0664)
	if err != nil {
		return err
	}

	return nil
}
