package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet, feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("unable to create request instance %v", err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("unable to make request %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("unable to read response body %v", err)
	}
	newRSSFeed := RSSFeed{}
	if err := xml.Unmarshal(body, &newRSSFeed); err != nil {
		return &RSSFeed{}, fmt.Errorf("unable to parse XML %v", err)
	}
	newRSSFeed.Channel.Title = html.UnescapeString(newRSSFeed.Channel.Title)
	newRSSFeed.Channel.Description = html.UnescapeString(newRSSFeed.Channel.Description)
	for _, rssfeed := range newRSSFeed.Channel.Item {
		rssfeed.Title = html.UnescapeString(rssfeed.Title)
		rssfeed.Description = html.UnescapeString(rssfeed.Description)
	}
	return &newRSSFeed, nil
}
