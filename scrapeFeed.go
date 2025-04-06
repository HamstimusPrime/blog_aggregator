package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hamstimusprime/blog_aggregator/internal/database"
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

	for _, itemRSS := range fetchedFeedRSS.Channel.Item {

		publishedAt, err := time.Parse(time.RFC1123Z, itemRSS.PubDate)
		if err != nil {
			publishedAt = time.Time{}
			fmt.Printf("could not parse date: %v, error: %v", itemRSS.PubDate, err)
		}
		createPostParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       itemRSS.Title,
			Url:         itemRSS.Link,
			Description: itemRSS.Description,
			PublishedAt: publishedAt,
			FeedID:      nextFeedID,
		}

		_, err = s.db.CreatePost(context.Background(), createPostParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post from feed %v with title %q: %v", nextFeedID, itemRSS.Title, err)
			continue
		}

	}
	return nil
}
