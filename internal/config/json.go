package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const gatorconfigFile = ".gatorconfig.json"

// gets the filepath of gatorconfigFile on the home directory
func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("your directories are homeless. %v", err)
	}
	filepath := home + "/" + gatorconfigFile
	return filepath, nil
}

// turns the json at the gatorconfigFile to a Config struct
func Read() (Config, error) {
	filepath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jason, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("~/%v couldn't be opened: %v", gatorconfigFile, err)
	}

	var config Config
	err = json.Unmarshal(jason, &config)
	if err != nil {
		return Config{}, fmt.Errorf("~/%v couldn't be unmarshaled: %v", gatorconfigFile, err)
	}
	return config, nil
}

// takes a Config struct and overwrites the gatorconfigFile with its data in JSON
func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	jason, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("honestly surprised you got this error. %v", err)
	}

	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, jason, 0644)
	return err
}
