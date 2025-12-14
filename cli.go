package main

import (
	"errors"
	"fmt"
	"github.com/trollian-alien/blog_aggregator/internal/config"
)

// internal state handler
type state struct{
	cfg *config.Config
}

// commands will take this shape
type command struct {
	name string
	args []string
}

// logins!!!
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("invalid command; no arguments")
	}
	username := cmd.args[0]
	err := s.cfg.SetUser(username)
	if err != nil {return err}
	fmt.Printf("User %v set\n", username)
	return nil
}

// struct of all commands
type commands struct {
	cmds map[string]func(*state, command) error
}

//command run method
func(c * commands) run(s *state, cmd command) error {
	commandFunction, ok := c.cmds[cmd.name]
	if !ok {
		return errors.New("this command does not exist")
	}
	return commandFunction(s, cmd)
}

//register new handling function for a command name
func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
