# nostr-bot

A multi-bot system that bridges RSS feeds to the Nostr network. Each bot independently monitors specific RSS feeds and republishes them as Nostr events.

## Features

- Multi-bot architecture with concurrent feed monitoring
- Each bot can have its own RSS feeds and Nostr private key
- YAML-based configuration for easy bot management
- Automatically syncs RSS feed items to the Nostr network
- Tracks published content using an SQLite database
- Prevents duplicate posts
- Graceful shutdown handling

## Requirements

- Go 1.x
- SQLite3

## Configuration

Create a `config/bots.yaml` file:
```yaml
relay_url: "wss://relay.damus.io"
database_path: "data/content.db"

bots:
  - name: "example-1"
    private_key: "nsec1..."
    rss_feeds:
      - "https://example.com/rss"

  - name: "example-2"
    private_key: "nsec2..."
    rss_feeds:
      - "https://example.com/rss"

```

The system will use `wss://relay.damus.io` as the default Nostr relay if none is specified.

## Setup & Running

1. Clone the repository
2. Create `config/bots.yaml` with your bot configurations
3. Run: `make run`
4. To clean the database: `make clean`

## Architecture

The system uses Go's concurrency features to manage multiple bots efficiently. Here's how it works:

```
Main Thread                     Bots goroutines
─────────────────               ──────────────────
│                              Bot 1
│ main()                       │ watchFeeds() -> RSS fetcher
│ ├─ Create bots               │
│ ├─ b.Start() ───────────────>│
│ ├─ b.Start() ───────────────>│ Bot 2
│ │                            │ watchFeeds() -> RSS fetcher
│ ├─ Create channels           │
│ │                           
│ └─ Block on <-c               
│     (wait for Ctrl+C)         
│                              
└────────────────              
```

### How it Works

1. **Initialization**
   - Main thread loads configuration and creates bot instances
   - Each bot is configured with its own private key and RSS feed list
   - Database connection is shared among all bots

2. **Concurrent Operation**
   - Each bot runs in its own goroutine
   - Bots independently monitor their RSS feeds
   - Main thread waits for shutdown signal

3. **Feed Processing**
   - Bots check their feeds every 5 minutes
   - New items are published to Nostr network
   - Published content is tracked in SQLite database
   - Duplicate posts are prevented

4. **Graceful Shutdown**
   - System catches interrupt signal (Ctrl+C)
   - All bots are notified to stop
   - Resources are properly cleaned up

## Project Structure

```
├── cmd/
│   └── nostr-bot/      # Application entry point
├── config/             # Configuration management
├── internal/           # Internal packages
│   ├── bot/           # Bot implementation
│   ├── database/      # Database operations
│   ├── nostr/         # Nostr client
│   └── rss/           # RSS feed handling
└── config/
    └── bots.yaml      # Bot configurations
```

## License
We don't care
