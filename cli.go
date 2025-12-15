package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/trollian-alien/blog_aggregator/internal/config"
	"github.com/trollian-alien/blog_aggregator/internal/database"
)

// internal state handler
type state struct{
	db  *database.Queries
	cfg *config.Config
}

// commands will take this shape
type command struct {
	name string
	args []string
}

// user login handler
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("invalid command; no arguments")
	}
	username := cmd.args[0]

	//check if user exists
	_, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		fmt.Println("User does not exist!")
		os.Exit(1)
	} else if err != nil {
		return fmt.Errorf("shrodinger's user. error: %v", err)
	}

	//setting user
	err = s.cfg.SetUser(username)
	if err != nil {return err}
	fmt.Printf("User %v set\n", username)
	return nil
}

// user registration handler
func handlerRegister(s *state, cmd command) error {
	//mandatory command check
	if len(cmd.args) == 0 {
		return errors.New("invalid command; no arguments")
	}
	username := cmd.args[0]

	//check if user exists
	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		fmt.Println("User already exists!")
		os.Exit(1)
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("shrodinger's user. error: %v", err)
	}

	//user creation
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name: username,
	})
	if err != nil {return err}
	fmt.Println("User created! User:")
	fmt.Printf("ID: %v\n", user.ID)
	fmt.Printf("Name: %v\n", user.Name)
	fmt.Printf("CreatedAt: %v\n", user.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", user.UpdatedAt)

	//setting user
	err = s.cfg.SetUser(username)
	if err != nil {return err}
	fmt.Printf("User %v set\n", username)
	return nil
}

//reset is a dangerous command that deletes all users. enjoy!
func handlerReset(s* state, cmd command) error {
	//this command ignores arguments
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("deletion problems?! %v", err)
	}
	fmt.Println("Reset Succesful!")
	return nil
}

//users command handler; gets all registered user names
func handlerUsers(s *state, cmd command) error {
	usernames, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("can't list users: %v", err)
	}
	for _, username := range usernames {
		if username == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", username)
		} else {
			fmt.Printf("* %v\n", username)
		}
	}
	return nil
}

//agg command handler
func handlerAgg(s *state, cmd command) error {
	 feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	 fmt.Println(feed)
	 return err
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
