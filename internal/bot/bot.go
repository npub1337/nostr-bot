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
		stopChan:    make(chan struct{}),
		rssFetcher:  rss.NewFetcher(),
	}, nil
}

func (b *Bot) Start() {
	go b.watchFeeds()
}

func (b *Bot) Stop() {
	close(b.stopChan)
}

func (b *Bot) watchFeeds() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.checkAndPublishUpdates()
		case <-b.stopChan:
			return
		}
	}
}

func (b *Bot) checkAndPublishUpdates() {
	items, err := b.rssFetcher.Fetch(b.RSSFeeds)
	if err != nil {
		log.Printf("[Bot: %s] Error fetching RSS: %v", b.Name, err)
		return
	}

	for _, item := range items {
		if b.DB.IsContentStored(item.ID) {
			continue
		}

		err := b.DB.InsertRetrievedContent(item.ID, item.Content, item.Source)
		if err != nil {
			log.Printf("[Bot: %s] Failed to insert content: %v", b.Name, err)
			continue
		}

		// Converter RSS Item para Nostr Content
		nostrContent := nostr.Content{
			ID:      item.ID,
			Content: item.Content,
			Source:  item.Source,
		}

		b.NostrClient.PublishContent(b.DB, []nostr.Content{nostrContent})
	}
}
