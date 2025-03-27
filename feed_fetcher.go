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
	//fetch feed from url
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
	/*the body of the response would need to be parsed into bytes using the
	io.ReadAll(res.body) function. it is the bytes that we would then pass to
	unmarshall alongside a pointer to an instance of the struct we want to convert into*/
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("unable to read response body %v", err)
	}
	newRSSFeed := RSSFeed{}
	if err := xml.Unmarshal(body, &newRSSFeed); err != nil {
		return &RSSFeed{}, fmt.Errorf("unable to parse XML %v", err)
	}
	/* by the time the body is parsed inot the rss feed struct, the title and description MIGHT contain
	tags like this '&ldquo;' this tag is an escape string, and we need to remove it.
	to do that, we use the UnescapeString of the html package to filter it out. we do this
	also with the items field*/
	newRSSFeed.Channel.Title = html.UnescapeString(newRSSFeed.Channel.Title)
	newRSSFeed.Channel.Description = html.UnescapeString(newRSSFeed.Channel.Description)
	for _, rssfeed := range newRSSFeed.Channel.Item {
		rssfeed.Title = html.UnescapeString(rssfeed.Title)
		rssfeed.Description = html.UnescapeString(rssfeed.Description)
	}
	return &newRSSFeed, nil
}
