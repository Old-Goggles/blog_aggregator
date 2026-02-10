package main

import (
	"context"
	"fmt"

	"github.com/Old-Goggles/blog_aggregator/internal/database"
)

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
