package main

import (
	"context"
	"fmt"
	"time"
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
		fmt.Println(item.Title)
	}

}
