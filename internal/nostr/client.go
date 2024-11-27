package nostr

import (
	"context"
	"fmt"
	"strings"
	"time"

	"nostr-bot/internal/database"

	"github.com/nbd-wtf/go-nostr"
)

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

func (client *Client) PublishContent(db *database.DB, content Content) error {
	event := nostr.Event{
		PubKey:    client.PublicKey,
		CreatedAt: nostr.Timestamp(time.Now().Unix()),
		Kind:      1,
		Content:   content.Content,
	}

	event.Sign(client.PrivateKey)

	err := client.Relay.Publish(context.Background(), event)
	if err != nil {
		if strings.Contains(err.Error(), "rate-limited") {
			return fmt.Errorf("rate-limited")
		}
		return err
	}

	return nil
}
