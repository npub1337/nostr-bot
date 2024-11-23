package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"xNostrBot/database"

	"github.com/joho/godotenv"
	"github.com/nbd-wtf/go-nostr"
)

type Config struct {
	TwitterBearerToken string
	NostrPrivateKey    string
	TwitterUserID      string
}

type Tweet struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type TwitterResponse struct {
	Data []Tweet `json:"data"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dbPath := filepath.Join("data", "tweets.db")
	os.MkdirAll(filepath.Dir(dbPath), 0755)

	db, err := database.New(dbPath)
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	config := Config{
		TwitterBearerToken: os.Getenv("TWITTER_BEARER_TOKEN"),
		NostrPrivateKey:    os.Getenv("NOSTR_PRIVATE_KEY"),
		TwitterUserID:      os.Getenv("TWITTER_USER_ID"),
	}

	nostrPrivKey := config.NostrPrivateKey
	nostrPubKey, err := nostr.GetPublicKey(nostrPrivKey)
	if err != nil {
		log.Fatal("Error getting Nostr public key:", err)
	}

	//coloquei o relay.damus.io, mas depois temos que por mais relays. Se eu não me
	// engano os relays repassam as mensagens para outros relays
	relay, err := nostr.RelayConnect(context.Background(), "wss://relay.damus.io")
	if err != nil {
		log.Fatal("Error connecting to Nostr relay:", err)
	}

	tweets, err := fetchTweets(config.TwitterBearerToken, config.TwitterUserID)
	if err != nil {
		log.Fatalf("Error fetching tweets: %v", err)
	}

	for _, tweet := range tweets {
		if db.IsTweetPublished(tweet.ID) {
			log.Printf("Tweet %s already published, skipping...", tweet.ID)
			continue
		}

		event := nostr.Event{
			PubKey:    nostrPubKey,
			CreatedAt: nostr.Timestamp(time.Now().Unix()),
			Kind:      1,
			Content:   tweet.Text,
		}

		event.Sign(nostrPrivKey)

		err = relay.Publish(context.Background(), event)
		if err != nil {
			log.Printf("Error publishing to Nostr: %v", err)
			continue
		}
		err = db.MarkTweetAsPublished(database.Tweet{
			ID:   tweet.ID,
			Text: tweet.Text,
		})
		if err != nil {
			log.Printf("Error marking tweet as published: %v", err)
			continue
		}

		log.Printf("Published tweet to Nostr: %s", tweet.Text)
		time.Sleep(time.Second)
	}
}

func fetchTweets(bearerToken, userID string) ([]Tweet, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitter.com/2/users/%s/tweets", userID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Twitter API returned status %s", resp.Status)
	}

	var twitterResponse TwitterResponse
	err = json.NewDecoder(resp.Body).Decode(&twitterResponse)
	if err != nil {
		return nil, err
	}

	return twitterResponse.Data, nil
}
