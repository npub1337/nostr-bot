package bot

import (
	"log"
	"time"

	"nostr-bot/internal/database"
	"nostr-bot/internal/nostr"
	"nostr-bot/internal/rss"
)

func NewBot(name string, privateKey string, relayURL string, feeds []string, db *database.DB) (*Bot, error) {
	nostrClient, err := nostr.NewClient(privateKey, relayURL)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Name:        name,
		NostrClient: nostrClient,
		RSSFeeds:    feeds,
		DB:          db,
		rssFetcher:  rss.NewFetcher(),
	}, nil
}

func (b *Bot) Start() {
	b.checkRSSFeeds()
	b.publishPendingItems()
}

func (b *Bot) checkRSSFeeds() {
	log.Printf("[Bot: %s] Starting RSS feeds check", b.Name)

	items, err := b.rssFetcher.Fetch(b.RSSFeeds)
	if err != nil {
		log.Printf("[Bot: %s] Error fetching RSS: %v", b.Name, err)
		return
	}

	for _, item := range items {
		if b.DB.IsContentStored(item.ID) {
			continue
		}

		log.Printf("[Bot: %s] New item found: %s", b.Name, item.ID)
		err := b.DB.InsertRetrievedContent(item.ID, item.Content, item.Source)
		if err != nil {
			log.Printf("[Bot: %s] Failed to insert content: %v", b.Name, err)
		}
	}
}

func (b *Bot) publishPendingItems() {
	pendingItems, err := b.DB.GetPendingContent()
	if err != nil {
		log.Printf("[Bot: %s] Error getting pending content: %v", b.Name, err)
		return
	}

	for _, item := range pendingItems {
		log.Printf("[Bot: %s] Attempting to publish: %s", b.Name, item.ID)

		nostrContent := nostr.Content{
			ID:      item.ID,
			Content: item.Content,
			Source:  item.Source,
		}

		err := b.NostrClient.PublishContent(b.DB, nostrContent)
		if err != nil {
			if err.Error() == "rate-limited" {
				log.Printf("[Bot: %s] Rate limited, will retry later: %s", b.Name, item.ID)
				return
			}
			log.Printf("[Bot: %s] Failed to publish: %v", b.Name, err)
			b.DB.UpdateContentStatus(item.ID, "failed")
			continue
		}

		err = b.DB.UpdateContentStatus(item.ID, "published")
		if err != nil {
			log.Printf("[Bot: %s] Error updating content status: %v", b.Name, err)
		}

		time.Sleep(1 * time.Second)
	}
}
