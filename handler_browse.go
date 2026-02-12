package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Old-Goggles/blog_aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command) error {
	ctx := context.Background()
	limit := 2
	if len(cmd.Args) > 0 {
		num, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("error setting limit %w", err)
		}
		limit = num
	}

	user, err := s.Db.GetUser(ctx, s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("unable to get user %w", err)
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	result, err := s.Db.GetPostsForUser(ctx, params)
	if err != nil {
		return fmt.Errorf("unable to get posts %w", err)
	}

	for _, post := range result {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Println("---")
	}

	return nil
}
