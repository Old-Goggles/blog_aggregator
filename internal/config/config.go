package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	fullPath := home + "/.gatorconfig.json"

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	fullPath := home + "/.gatorconfig.json"

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(fullPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
