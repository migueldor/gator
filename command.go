package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commandHandlers map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Username cannot be empty\n")
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("Logged as: %s\n", cmd.args[0])
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	if handler, ok := c.commandHandlers[cmd.name]; !ok {
		return fmt.Errorf("%s command doesn't exists", cmd.name)
	} else {
		return handler(s, cmd)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandHandlers[name] = f
}
