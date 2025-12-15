package main

import (
	"context"
	"encoding/xml"
	"html"
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
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {return &RSSFeed{}, err}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {return &RSSFeed{}, err}
	defer res.Body.Close()

	feed := &RSSFeed{}
	decoder := xml.NewDecoder(res.Body)
	err = decoder.Decode(feed)
	if err != nil {return &RSSFeed{}, err}

	//handle escpaed html entities
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}

	return feed, nil
}