package main

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

func (feed *RSSFeed) unEscape()  {

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title) 
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description) 

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
}

func (feed *RSSFeed) print() {

	fmt.Println()
	fmt.Printf("Title: %s\n", feed.Channel.Title)
	fmt.Println()
	fmt.Printf("Link: %s\n", feed.Channel.Link)
	fmt.Println()
	fmt.Printf("Description: %s\n", feed.Channel.Description)
	fmt.Println()

	for i, v := range feed.Channel.Item {
		fmt.Printf("Item %d:\n", i)
		fmt.Printf(" - Title: %s\n", v.Title)
		fmt.Printf(" - Link: %s\n", v.Link)
		fmt.Printf(" - Description: %s\n", v.Description)
		fmt.Printf(" - Publication Date: %s\n", v.PubDate)
		fmt.Println()
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}	
	req.Header.Set("User-Agent", "gator")
	
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {	
		return &RSSFeed{}, err
	}
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return &RSSFeed{}, err
	}
	
	
	return &feed, nil
}
