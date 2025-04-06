package main

import (
	"errors"
	"fmt"
)

func containArgs(cmd command) error {
	if len(cmd.Args) == 0 {
		fmt.Printf("%s command expects args. No arguments passed", cmd.Name)
		return errors.New("no args passed")
	}
	return nil
}
