package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/Ahmed0427/rssy/internal/database"
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

func parseRSSTimeFromat(t string) time.Time {
	timeLayouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
	}

	for _, layout := range timeLayouts {
		parsedTime, err := time.Parse(layout, t)
		if err == nil {
			return parsedTime
		}
	}

	return time.Now()
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "rssy")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	feed := &RSSFeed{}
	bodyData, err := io.ReadAll(res.Body)
	err = xml.Unmarshal(bodyData, feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := 0; i < len(feed.Channel.Item); i++ {
		feed.Channel.Item[i].Title =
			html.UnescapeString(feed.Channel.Item[i].Title)

		feed.Channel.Item[i].Description =
			html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return feed, nil
}

func scrapeFeeds(ctx context.Context, s *State) (*RSSFeed, *database.Feed, error) {
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return nil, nil, err
	}

	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return nil, nil, err
	}

	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Url: feed.Url,
	})
	if err != nil {
		return nil, nil, err
	}

	return rssFeed, &feed, nil
}
