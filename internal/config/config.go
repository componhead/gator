package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	CurrentUserName string `json:"current_user_name"`
	DbUrl           string `json:"db_url"`
}

func Read() (Config, error) {
	cfg := Config{}
	p, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	f, err := os.ReadFile(p)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(f, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func write(cfg Config) error {
	p, err := getConfigFilePath()
	if err != nil {
		return err
	}

	s, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(p, s, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir = fmt.Sprint(dir, string(os.PathSeparator), configFileName)
	return dir, nil
}

func (cfg *Config) SetUser(usr string) error {
	cfg.CurrentUserName = usr
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}
