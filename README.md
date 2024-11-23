# xNostrBot

A bot that bridges Twitter/X posts to the Nostr network. It automatically fetches tweets from a specified Twitter account and republishes them as Nostr events.

## Features

- Automatically syncs tweets to the Nostr network
- Tracks published tweets using an SQLite database
- Prevents duplicate posts
- Easy setup using environment variables

## Requirements

- Go 1.x
- SQLite3
- Twitter/X API Bearer Token
- Nostr Private Key

## Environment Variables

Create a .env file with the following variables:

```
TWITTER_BEARER_TOKEN=your_twitter_bearer_token  
NOSTR_PRIVATE_KEY=your_nostr_private_key  
TWITTER_USER_ID=twitter_user_id_to_monitor  
```

## Setup & Running

1. Clone the repository:  
2. Create a .env file with the required variables.
3. Do: make run 
4. To clean the database use: make clean

## How it Works

1. The bot connects to Twitter's API and fetches recent tweets from the specified user.
2. Each new tweet is published to the Nostr network using the provided private key.
3. Published tweets are tracked in a local SQLite database to prevent duplicates.
4. Currently configured to publish to relay.damus.io.

## License
 I don't care
