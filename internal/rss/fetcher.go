package rss

import (
	"github.com/mmcdole/gofeed"
)

type RSSFetcher struct{}

func NewFetcher() *RSSFetcher {
	return &RSSFetcher{}
}

func (f *RSSFetcher) Fetch(feedURLs []string) ([]Item, error) {
	var items []Item
	parser := gofeed.NewParser()

	for _, url := range feedURLs {
		feed, err := parser.ParseURL(url)
		if err != nil {
			continue // Podemos logar o erro aqui se necessário
		}

		for _, entry := range feed.Items {
			items = append(items, Item{
				ID:      entry.GUID,
				Title:   entry.Title,
				Link:    entry.Link,
				Content: entry.Title + "\n" + entry.Link,
				Source:  "rss",
			})
		}
	}

	return items, nil
}
