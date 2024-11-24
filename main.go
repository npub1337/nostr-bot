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
	TwitterBearerToken string
	NostrPrivateKey    string
	TwitterUserID      string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dbPath := filepath.Join("data", "content.db")
	os.MkdirAll(filepath.Dir(dbPath), 0755)

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	// TODO: add configs (e.g. URLs, keys)
	config := Config{
		NostrPrivateKey: os.Getenv("NOSTR_PRIVATE_KEY"),
	}

	/* nostr */

	nostrPrivKey := config.NostrPrivateKey
	nostrCli, err := nostr.NewClient(nostrPrivKey)
	if err != nil {
		log.Fatal("Error getting Nostr public key:", err)
	}

	/* fetchers */

	// TODO: add other sources
	rssFeeds := []string{
		"https://feedx.net/rss/ap.xml",
		// "https://feeds.bbci.co.uk/news/world/rss.xml",
		// "https://www.euronews.com/rss",
		// "https://www.lemonde.fr/en/rss/une.xml",
		// "https://time.com/feed/"
	}

	rssItems, err := rss.FetchRSSFeeds(rssFeeds)
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

		err := db.InsertRetrievedContent(contentID, item.Content, "rss")
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

		log.Printf("Published content: %s", item)
	}

}

func generateContentID(content string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(content)))
}
