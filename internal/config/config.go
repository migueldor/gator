package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, jsonData, 0644)
	//the code 0644 is a permition code, this code is the standard
	if err != nil {
		return err
	}
	return nil

}
