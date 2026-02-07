package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Old-Goggles/blog_aggregator/internal/config"
	"github.com/Old-Goggles/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	Db  *database.Queries
	Cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		username := s.Cfg.CurrentUserName
		user, err := s.Db.GetUser(ctx, username)
		if err != nil {
			return fmt.Errorf("error finding user in database %w", err)
		}
		return handler(s, cmd, user)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is required")
	}

	ctx := context.Background()
	username := cmd.Args[0]

	_, err := s.Db.GetUser(ctx, username)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user name does not exist")
	} else if err != nil {
		return fmt.Errorf("error finding user in database %w", err)
	}

	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Println("user has been set to", username)
	return nil
}

type commands struct {
	Handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found")
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.Handlers == nil {
		c.Handlers = make(map[string]func(*state, command) error)
	}
	c.Handlers[name] = f
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is required")
	}

	username := cmd.Args[0]
	ctx := context.Background()
	id := uuid.New()
	now := time.Now()
	params := database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      username,
	}

	result, err := s.Db.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("user name already exists: %w", err)
	}

	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("User Created: %+v", result)
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.Db.DeleteUsers(ctx)
	if err != nil {
		return fmt.Errorf("unable to delete: %w", err)
	}

	fmt.Printf("users reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()

	result, err := s.Db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("unable to get users")
	}

	for _, user := range result {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()

	result, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("unable to fetch feed %v", err)
	}

	fmt.Printf("Feed: %+v\n", result)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	if len(cmd.Args) != 2 {
		return fmt.Errorf("name and url are required")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	params := database.CreateFeedParams{
		Name:   name,
		Url:    url,
		UserID: user.ID,
	}

	feed, err := s.Db.CreateFeed(ctx, params)
	if err != nil {
		return fmt.Errorf("error creating feed %w", err)
	}

	follow_params := database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	_, err = s.Db.CreateFeedFollow(ctx, follow_params)
	if err != nil {
		return fmt.Errorf("error creating feed follow")
	}

	fmt.Printf("feed created: %+v\n", feed)
	fmt.Printf("automatically following as %s\n", user.Name)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.Db.GetAllFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error getting feeds %w", err)
	}

	for _, feed := range feeds {
		user, err := s.Db.GetUserByID(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting user by ID %w", err)
		}
		fmt.Printf("Name: %s\nURL: %s\nUser: %s\n\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	if len(cmd.Args) != 1 {
		return fmt.Errorf("url is required")
	}

	url := cmd.Args[0]
	feed, err := s.Db.GetFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("error getting feed %w", err)
	}
	params := database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	followed, err := s.Db.CreateFeedFollow(ctx, params)
	if err != nil {
		return fmt.Errorf("error following feed %w", err)
	}
	fmt.Printf("%s is now following %s\n", followed.UserName, followed.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	following, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error getting followed feeds %w", err)
	}

	for _, feed := range following {
		fmt.Printf("%s\n", feed.FeedName)
	}
	return nil
}
