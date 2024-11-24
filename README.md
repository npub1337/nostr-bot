# nostr-bot

A bot that bridges RSS feeds to the Nostr network. It automatically fetches news from specified RSS feeds and republishes them as Nostr events.

## Features

- Automatically syncs RSS feed items to the Nostr network
- Tracks published content using an SQLite database
- Prevents duplicate posts
- Easy setup using environment variables

## Requirements

- Go 1.x
- SQLite3
- Nostr Private Key

## Environment Variables

Create a .env file with the following variables:

```
NOSTR_PRIVATE_KEY=your_nostr_private_key
NOSTR_RELAY_URL=your_nostr_relay_url
```

## Setup & Running

1. Clone the repository
2. Create a .env file with the required variables
3. Do: make run 
4. To clean the database use: make clean

## How it Works

1. The bot connects to configured RSS feeds and fetches recent news items
2. Each new item is published to the Nostr network using the provided private key
3. Published content is tracked in a local SQLite database to prevent duplicates
4. Publishes to the configured relay (defaults to relay.damus.io if not specified)

## License
 I don't care
