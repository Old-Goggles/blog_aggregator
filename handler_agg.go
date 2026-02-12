package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Old-Goggles/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("duration is required")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("unable to parse duration %w", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)
	fmt.Printf("collecting feeds every %v", timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	ctx := context.Background()

	feed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Printf("unable to get next feed %v", err)
		return
	}

	err = s.Db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		fmt.Printf("unable to mark feed fetched %v", err)
		return
	}

	rss, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		fmt.Printf("unable to fetch feed: %v", err)
		return
	}

	fmt.Printf("Found feed: %s\n", feed.Name)
	for _, item := range rss.Channel.Item {

		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt.Time = t
			publishedAt.Valid = true
		}

		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		}

		_, err = s.Db.CreatePost(ctx, params)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				continue
			}
			fmt.Printf("Error creating post: %v\n", err)
			continue
		}
	}

}
