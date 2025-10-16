package main

import (
	"errors"
	"fmt"

	"github.com/componhead/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	var username string
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	username = cmd.Args[0]
	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("Username %s setted\n", username)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
