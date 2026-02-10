package main

import (
	"context"
	"fmt"

	"github.com/Old-Goggles/blog_aggregator/internal/database"
)

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

func handlerUnfollow(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	if len(cmd.Args) != 1 {
		return fmt.Errorf("url is required")
	}

	url := cmd.Args[0]
	feed, err := s.Db.GetFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("error getting feed %w", err)
	}

	params := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err := s.Db.UnfollowFeed(ctx, params); err != nil {
		return fmt.Errorf("error unfollowing feed %w", err)
	}

	return nil
}
