package main

import (
	"fmt"
	"os"

	"github.com/hamstimusprime/blog_aggregator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config file: %v", err)
		return
	}

	s := state{
		config: &cfg,
	}

	cmds := commands{
		handlersMap: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	commandLineInput := os.Args
	if len(commandLineInput) < 3 {
		fmt.Println("error, no arguments provided")
		os.Exit(1)
	}
	if len(commandLineInput) > 3 {
		fmt.Println("warning! multiple arguments provided, only first argument would be used")
	}
	commandName := commandLineInput[1]
	commandArgs := commandLineInput[2:]
	cmd := command{name: commandName, args: commandArgs}

	if err = cmds.run(&s, cmd); err != nil {
		fmt.Println("could not execute run command")
		os.Exit(1)
	}
}

type state struct {
	config *config.Config
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("expected arguments. No arguments provided")
	}

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		fmt.Printf("unable to set user name, %v", err)
		return err
	}
	fmt.Println("username has been set")
	return nil
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlersMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlersMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlersMap[cmd.name]
	if !ok {
		return fmt.Errorf("invalid command")
	}
	handler(s, cmd)
	return nil
}
