package nostr

import (
	"context"
	"log"
	"time"

	"nostr-bot/database"
	"nostr-bot/pkg/models"

	"github.com/nbd-wtf/go-nostr"
)

type Client struct {
	PrivateKey string
	PublicKey  string
	Relay      *nostr.Relay
}

func NewClient(privateKey, relayURL string) (*Client, error) {
	pubKey, err := nostr.GetPublicKey(privateKey)
	if err != nil {
		return nil, err
	}

	relay, err := nostr.RelayConnect(context.Background(), relayURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		PrivateKey: privateKey,
		PublicKey:  pubKey,
		Relay:      relay,
	}, nil
}

func (client *Client) PublishContent(db *database.DB, content []models.Content) {
	for _, item := range content {
		if db.IsContentStored(item.ID) {
			log.Printf("Content %s already published, skipping...", item.ID)
			continue
		}

		event := nostr.Event{
			PubKey:    client.PublicKey,
			CreatedAt: nostr.Timestamp(time.Now().Unix()),
			Kind:      1,
			Content:   item.Content,
		}

		event.Sign(client.PrivateKey)

		err := client.Relay.Publish(context.Background(), event)
		if err != nil {
			log.Printf("Error publishing to Nostr: %v", err)
			continue
		}

		err = db.MarkAsPublished(item.Content)
		if err != nil {
			log.Printf("Error setting Publish status on DB: %v", err)
			continue
		}

		log.Printf("Published content to Nostr: %s", item.Content)
		time.Sleep(time.Second)
	}
}
