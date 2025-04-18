package feeds

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// declaring request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)

	// checking for error
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// This request is officially a gator, RAAAAAHHH! 🐊
	req.Header.Set("User-Agent", "gator")

	// declaring default client
	client := &http.Client{}

	// doing request -> getting response
	resp, err := client.Do(req)

	// checking for error
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}

	// cleaning my mess (WHO MADE THAT MESS?! -> you did, king 👑)
	defer resp.Body.Close()

	// check for status for safety
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// reading the response body
	body, err := io.ReadAll(resp.Body)

	// checking for error
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// creating a new RSSFeed struct instance
	var my_feed RSSFeed

	// unmarshalling the XML response into the struct instance
	err = xml.Unmarshal(body, &my_feed)

	// checking for error
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %w", err)
	}

	// clean the channel info
	my_feed.Channel.Title = html.UnescapeString(my_feed.Channel.Title)
	my_feed.Channel.Description = html.UnescapeString(my_feed.Channel.Description)

	// clean the items info
	for i := range my_feed.Channel.Item {
		my_feed.Channel.Item[i].Title = html.UnescapeString(my_feed.Channel.Item[i].Title)
		my_feed.Channel.Item[i].Description = html.UnescapeString(my_feed.Channel.Item[i].Description)
	}

	// returning clean feed
	return &my_feed, nil
}


