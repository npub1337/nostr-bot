package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"nostr-bot/database"
	"nostr-bot/pkg/models"
	"nostr-bot/pkg/nostr"
	"nostr-bot/pkg/rss"

	"github.com/joho/godotenv"
)

type Config struct {
	NostrPrivateKey string
	NostrRelayURL   string
	RSSFeeds        []string
	DatabasePath    string
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		NostrPrivateKey: os.Getenv("NOSTR_PRIVATE_KEY"),
		NostrRelayURL:   os.Getenv("NOSTR_RELAY_URL"),
		DatabasePath:    filepath.Join("data", "content.db"),
		RSSFeeds: []string{
			"https://feedx.net/rss/ap.xml",
			"https://feeds.bbci.co.uk/news/world/rss.xml",
			// TODO: add more feeds
		},
	}

	if config.NostrPrivateKey == "" {
		return nil, fmt.Errorf("NOSTR_PRIVATE_KEY is required")
	}

	if config.NostrRelayURL == "" {
		config.NostrRelayURL = "wss://relay.damus.io"
	}

	return config, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	os.MkdirAll(filepath.Dir(config.DatabasePath), 0755)

	db, err := database.InitDB(config.DatabasePath)
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	nostrCli, err := nostr.NewClient(config.NostrPrivateKey, config.NostrRelayURL)
	if err != nil {
		log.Fatal("Error initializing Nostr client:", err)
	}

	rssItems, err := rss.FetchRSSFeeds(config.RSSFeeds)
	if err != nil {
		log.Fatal("Error fetching RSS..:", err)
	}

	/* storing content */

	// TODO: add other sources
	allContent := append(rssItems)
	for _, item := range allContent {
		contentID := generateContentID(item.Content)

		if db.IsContentStored(contentID) {
			log.Printf("Skipping already stored content: %s", item)
			continue
		}

		err := db.InsertRetrievedContent(contentID, item.Content, item.Source)
		if err != nil {
			log.Printf("Failed to insert content into DB: %v", err)
			continue
		}
	}

	/* publisher */

	for _, item := range allContent {
		// TODO: build full content (i.e. remove loop)
		contents := []models.Content{
			{ID: item.ID, Content: item.Content},
		}

		nostrCli.PublishContent(db, contents)
		// if err != nil {
		// 	log.Fatalf("Failed to publish content: %v", err)
		// }
	}

}

func generateContentID(content string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(content)))
}
