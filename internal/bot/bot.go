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
	log.Printf("[Bot: %s] Bot started, will check feeds every 30 seconds", b.Name)
	go b.watchFeeds()
	go b.checkAndPublishUpdates()
}

func (b *Bot) Stop() {
	close(b.stopChan)
}

func (b *Bot) watchFeeds() {
	ticker := time.NewTicker(30 * time.Second)
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
	log.Printf("[Bot: %s] Starting RSS feeds check", b.Name)

	items, err := b.rssFetcher.Fetch(b.RSSFeeds)
	if err != nil {
		log.Printf("[Bot: %s] Error fetching RSS: %v", b.Name, err)
		return
	}
	log.Printf("[Bot: %s] Found %d items in RSS feeds", b.Name, len(items))

	for _, item := range items {
		if b.DB.IsContentStored(item.ID) {
			log.Printf("[Bot: %s] Item already exists in database: %s", b.Name, item.ID)
			continue
		}

		log.Printf("[Bot: %s] New item found: %s", b.Name, item.ID)
		err := b.DB.InsertRetrievedContent(item.ID, item.Content, item.Source)
		if err != nil {
			log.Printf("[Bot: %s] Failed to insert content: %v", b.Name, err)
			continue
		}

		nostrContent := nostr.Content{
			ID:      item.ID,
			Content: item.Content,
			Source:  item.Source,
		}

		log.Printf("[Bot: %s] Publishing content to Nostr: %s", b.Name, item.ID)
		b.NostrClient.PublishContent(b.DB, []nostr.Content{nostrContent})
	}
}
