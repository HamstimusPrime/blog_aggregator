package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const config_filename = "/.gatorconfig.json"

type Config struct {
	DbURL             string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.Current_user_name = username
	write(*c)
	return nil

}

func Read() (Config, error) {
	config_file_path, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("error getting config file path: %v", err)
		return Config{}, err
	}

	file, err := os.Open(config_file_path)
	if err != nil {
		fmt.Printf("error opening config file: %v", config_file_path)
		return Config{}, err
	}
	defer file.Close()

	config_struct := Config{}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config_struct); err != nil {
		fmt.Printf("error decoding request body, error: %v", err)
		return Config{}, err
	}
	return config_struct, nil
}

func getConfigFilePath() (string, error) {
	home_path, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("invalid home path")
		return "", fmt.Errorf("home path error")
	}
	config_file_path := home_path + config_filename
	return config_file_path, nil
}

func write(cfg Config) error {

	//open the file first
	config_file_path, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("error geting config file path, %v", err)
		return err
	}
	file, err := os.OpenFile(config_file_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
