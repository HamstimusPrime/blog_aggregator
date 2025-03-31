package main

import "fmt"

func containArgs(cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("command expects arguments. no arguments passed")
	}
	return nil
}
