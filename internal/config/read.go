package config

import (
	"encoding/json"
	"os"
)

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := home + "/.gatorconfig.json"
	return path, nil
}

func Read() (Config, error) {
	var newConfig Config
	path, err := getConfigFilePath()
	if err != nil {
		return newConfig, err
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return newConfig, err
	}
	err = json.Unmarshal(file, &newConfig)
	if err != nil {
		return newConfig, err
	}
	return newConfig, nil
}
