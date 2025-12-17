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
	//mandatory command check
	if len(cmd.args) == 0 {
		return errors.New("invalid command; no arguments")
	}
	timeInterval, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		fmt.Printf("Please enter a valid time duaration")
		fmt.Printf("for instance 10s or 2m or 3h")
		os.Exit(1)
	}

	fmt.Println("Collecting feeds every " + cmd.args[0])
	for range time.Tick(timeInterval) {
		err = scrapeFeeds(s)
		if err != nil {
			fmt.Println("Error encountered. Stopping now")
			return err
		}
		fmt.Println("OK")
	}
	return nil
}

//addfeed command handler
func handlerAddfeed(s *state, cmd command) error {
	//mandatory command check
	if len(cmd.args) != 2 {
		return errors.New("exactly two arguments expected, name and url respectively")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	username := s.cfg.CurrentUserName
	userID, err := s.db.GetUserID(context.Background(), username)
	if err != nil {
		return fmt.Errorf("couldn't get your user ID. Error: %v", err)
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID        : uuid.New(),
		CreatedAt : time.Now().UTC(),
		UpdatedAt : time.Now().UTC(),
		Name      : feedName,
		Url       : feedURL,
		UserID    : userID,
	})
	if err != nil {return err}

	fmt.Println("RSS feed created! Feed:")
	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("URL: %v\n", feed.Url)
	fmt.Printf("CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", feed.UpdatedAt)

	err = followFeed(s, feedURL, userID)
	return err
}

//feeds command handler
func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("can't fetch the feeds. Error: %v", err)
	}
	fmt.Println("Feeds:")
	for _, feed := range feeds {
		fmt.Println(feed)
	}
	return nil
}

// helper for follow and addfeed commands
func followFeed(s *state, feedURL string, userID uuid.UUID) error {
	feedID, err := s.db.FindFeedID(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error finding feed ID. %v", err)
	}
	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID       : uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID   : userID,
		FeedID   : feedID,
	})
	if err != nil {
		return fmt.Errorf("error following. %v", err)
	}
	
	fmt.Printf("You, %v, are now following %v!\n", follow.Username, follow.FeedName)
	return nil
}

//follow command handler
func handlerFollow(s* state, cmd command) error {
	if len(cmd.args) == 0 {
		fmt.Println("Please provide a feed url")
		os.Exit(1)
	}
	feedURL := cmd.args[0]
	userID, err := s.db.GetUserID(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error finding your user ID. %v", err)
	}
	err = followFeed(s, feedURL, userID)
	return err
}

//unfollow command handler
func handlerUnfollow(s* state, cmd command) error {
if len(cmd.args) == 0 {
		fmt.Println("Please provide a feed url")
		os.Exit(1)
	}
	feedURL := cmd.args[0]

	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		Name: s.cfg.CurrentUserName,
		Url: feedURL,
	})
	if err != nil {
		return fmt.Errorf("looks like you're stuck following that feed! Error: %v", err)
	}
	fmt.Println("Succesfully unfollowed feed!")
	return nil
}

//following command handler
func handlerFollowing(s* state, cmd command) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("trouble getting your followed feeds. %v", err)
	}
	fmt.Println("Your followed feeds:")
	for _, f := range feeds {
		fmt.Println(f.FeedName)
	}
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
