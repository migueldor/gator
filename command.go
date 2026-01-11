package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/migueldor/gator/internal/database"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commandHandlers map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: login <name>\n")
	}
	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w\n", err)
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("Logged as: %s\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: register <name>\n")
	}
	newUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	newUser, err := s.db.CreateUser(context.Background(), newUserParams)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}
	if err := s.config.SetUser(newUserParams.Name); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Println("User created successfully:")
	fmt.Printf(" * ID:   %v\n", newUser.ID)
	fmt.Printf(" * Name: %v\n", newUser.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't reset users: %w", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't retrieve users: %w", err)
	}
	for i := 0; i < len(users); i++ {
		if users[i] == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", users[i])
		} else {
			fmt.Printf("* %s\n", users[i])
		}

	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't retrieve feed: %w", err)
	}
	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>\n")
	}
	currentUserName := s.config.CurrentUserName
	user, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return fmt.Errorf("couldn't retrieve user: %w", err)
	}
	newFeedParams := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}
	newFeed, err := s.db.AddFeed(context.Background(), newFeedParams)
	if err != nil {
		return fmt.Errorf("couldn't add feed: %w", err)
	}

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    newFeed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return fmt.Errorf("couldn't follow %s feed: %w", newFeed.Name, err)
	}

	fmt.Println("Feed added successfully:")
	fmt.Printf(" * ID:   %v\n", newFeed.ID)
	fmt.Printf(" * Name: %v\n", newFeed.Name)
	fmt.Printf(" * Url: %v\n", newFeed.Url)
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't retrieve feeds: %w", err)
	}
	for i := range feeds {
		fmt.Println(feeds[i])
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: follow <url>\n")
	}
	currentUserName := s.config.CurrentUserName
	url := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return fmt.Errorf("couldn't retrieve user: %w", err)
	}
	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't retrieve feed: %w", err)
	}
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	newFeedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("couldn't follow %s feed: %w", feed.Name, err)
	}
	fmt.Println("Feed followed successfully:")
	fmt.Printf(" * Feed name: %v\n", newFeedFollow.FeedName)
	fmt.Printf(" * User: %v\n", newFeedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	currentUserName := s.config.CurrentUserName
	user, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return fmt.Errorf("couldn't retrieve user: %w", err)
	}
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't retrieve followed feeds: %w", err)
	}
	for i := range follows {
		fmt.Printf("* %s\n", follows[i].FeedName)
	}
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
