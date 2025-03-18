package main

import (
	"fmt"

	"github.com/hamstimusprime/blog_aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config file: %v", err)
		return
	}

	if err := cfg.SetUser("lane"); err != nil {
		fmt.Printf("error setting user: %v", err)
	}

	updated_cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config file: %v", err)
		return
	}
	fmt.Println(updated_cfg.DbURL)
	fmt.Println(updated_cfg.Current_user_name)

}
