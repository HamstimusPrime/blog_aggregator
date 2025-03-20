package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFilename = "/.gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"CurrentUserName"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(c)

}

func write(cfg *Config) error {

	//open the file first
	configFilePath, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("error geting config file path, %v", err)
		return err
	}
	file, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return err
	}
	defer file.Close()

	//convert the struct into a Json and write
	err = json.NewEncoder(file).Encode(cfg)
	if err != nil {
		fmt.Printf("failed to write json to file, %v", err)
		return err
	}
	return nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("error getting config file path: %v", err)
		return Config{}, err
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Printf("error opening config file: %v", configFilePath)
		return Config{}, err
	}
	defer file.Close()

	configStruct := Config{}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&configStruct); err != nil {
		fmt.Printf("error decoding request body, error: %v", err)
		return Config{}, err
	}
	return configStruct, nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("invalid home path")
		return "", fmt.Errorf("home path error")
	}
	configFilePath := homePath + configFilename
	return configFilePath, nil
}
