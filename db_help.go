package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/colfarl/gator/internal/database"
	"github.com/google/uuid"
)

func printUser(u database.User) {
	fmt.Printf("ID: %v\n", u.ID)
	fmt.Printf("Time Created: %v\n", u.CreatedAt)
	fmt.Printf("Time Updated: %v\n", u.UpdatedAt)
	fmt.Printf("Name: %v\n", u.Name.String)
}

func prettyPost(p database.Post) {
	fmt.Printf("Title: %v\n", p.Title)
	fmt.Printf("Description: %s\n", p.Description.String)
	fmt.Printf("Link: %s\n", p.Url)
	fmt.Printf("Published: %v\n", p.PublishedAt)
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {	

	return func(s *state, cmd command) error {	

		currUser := s.CurrentState.CurrentUserName
		currUserInfo, err := s.db.GetUser(context.Background(), sql.NullString{String: currUser, Valid: currUser != ""})
		if err != nil {
			return err
		}	

		return handler(s, cmd, currUserInfo)
	}
}

func parseTimeAnyLayout(timeStr string) (time.Time, error){
	layouts := []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.Kitchen,
		time.Stamp,
		time.StampMicro,
		time.StampMilli,
		time.StampNano,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
	}

	for _, v := range layouts {
		time, err := time.Parse(v, timeStr)
		if err == nil {
			return time, nil
		}
	}

	var t time.Time
	return t, fmt.Errorf("no matching layout for: %s", timeStr)
}

func scrapeFeeds(s *state) error {
	
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	
	rssFeed, err := fetchFeed(context.Background(), nextFeed.Url.String)
	if err != nil {
		return err
	}
	
	err = s.db.MarkedFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return err
	}
	
	rssFeed.unEscape()
	for _, v := range rssFeed.Channel.Item {
		publishTime, err := parseTimeAnyLayout(v.PubDate)
		if err != nil {
			return err
		}

		params := database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: v.Title,
			Url: v.Link,
			Description: sql.NullString{String: v.Description, Valid: v.Description != ""},
			PublishedAt: publishTime,
			FeedID: nextFeed.ID,
		}

		_, err = s.db.CreatePost(context.Background(), params)
		if err != nil {
			return err
		}
	}

	return nil
}
