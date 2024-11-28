package rss

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/mmcdole/gofeed"
)

type RSSFetcher struct{}

func NewFetcher() *RSSFetcher {
	return &RSSFetcher{}
}

func (f *RSSFetcher) Fetch(name string, feedURLs []string) ([]Item, error) {
	var items []Item
	parser := gofeed.NewParser()

	for _, url := range feedURLs {
		feed, err := parser.ParseURL(url)
		if err != nil {
			continue
		}

		for _, entry := range feed.Items {
			hash := md5.Sum([]byte(entry.Title + entry.Link))

			items = append(items, Item{
				ID:      hex.EncodeToString(hash[:]),
				BotName: name,
				Title:   entry.Title,
				Link:    entry.Link,
				Content: entry.Title + "\n" + entry.Link,
				Source:  "rss",
			})
		}
	}

	return items, nil
}
