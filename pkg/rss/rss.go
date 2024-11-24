package rss

import (
	"log"
	"nostr-bot/pkg/models"

	"github.com/mmcdole/gofeed"
)

func FetchRSSFeeds(feedURLs []string) ([]models.Content, error) {
	var items []models.Content
	parser := gofeed.NewParser()

	for _, url := range feedURLs {
		feed, _ := parser.ParseURL(url)

		for _, entry := range feed.Items {
			items = append(items, models.Content{
				ID:      entry.GUID,
				Content: entry.Title + "\n" + entry.Link,
			})
			log.Printf("%v", items)
		}
	}

	return items, nil
}
