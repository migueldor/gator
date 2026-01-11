package main

import (
	"context"
	"encoding/xml"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var body io.Reader = nil
	var newFeed RSSFeed
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, body)
	if err != nil {
		return &newFeed, err
	}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &newFeed, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &newFeed, err
	}
	err = xml.Unmarshal(data, &newFeed)
	if err != nil {
		return &newFeed, err
	}
	newFeed.Channel.Title = html.UnescapeString(newFeed.Channel.Title)
	newFeed.Channel.Description = html.UnescapeString(newFeed.Channel.Description)

	for i := range newFeed.Channel.Item {
		newFeed.Channel.Item[i].Title = html.UnescapeString(newFeed.Channel.Item[i].Title)
		newFeed.Channel.Item[i].Description = html.UnescapeString(newFeed.Channel.Item[i].Description)
	}
	return &newFeed, nil

}
