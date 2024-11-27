package bot

import (
	"nostr-bot/internal/database"
	"nostr-bot/internal/nostr"
	"nostr-bot/internal/rss"
)

type Bot struct {
	Name        string
	NostrClient *nostr.Client
	RSSFeeds    []string
	DB          *database.DB
	rssFetcher  rss.Fetcher
}
