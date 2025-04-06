package main

import (
	"context"
	"fmt"
)

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("unaple to fetch next feed")
		return err
	}
	nextFeedID := nextFeed.ID
	markedFeed, err := s.db.MarkFeedFetched(context.Background(), nextFeedID)
	if err != nil {
		fmt.Println("unable to mark fetched feed")
		return err
	}
	url := markedFeed.Url
	fetchedFeedRSS, err := fetchFeed(context.Background(), url)
	if err != nil {
		fmt.Println("unable to get feed by url")
		return err
	}
	for _, item := range fetchedFeedRSS.Channel.Item {
		fmt.Printf("feed title: %v\n", item.Title)
	}
	return nil
}
