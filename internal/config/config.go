//Package config to manage credentials and user
package config

import (
	"os"
	"encoding/json"
)

// ======== Helpers ======== 
const configFileName = ".gatorconfig.json" 

func getConfigPath() (string, error) {
	result, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return result + "/" + configFileName, nil
}

func write(cfg Config) error {

	configFilePath, err := getConfigPath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	
	err = os.WriteFile(configFilePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
// ======== Exports ======== 

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}


func Read() (Config, error) {

	configFilePath, err := getConfigPath() 

	if err != nil {
		return Config{}, err
	}
	
	body, err := os.ReadFile(configFilePath) 
	if err != nil {
		return Config{}, err
	}
	
	var config Config
	err = json.Unmarshal(body, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}


func (cfg *Config) SetUser(newUser string) error {		

	cfg.CurrentUserName = newUser

	err := write(*cfg)
	if err != nil {
		return err
	}

	return nil
}


