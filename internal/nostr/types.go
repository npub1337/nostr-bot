package nostr

import (
	"github.com/nbd-wtf/go-nostr"
)

type Client struct {
	PrivateKey string
	PublicKey  string
	Relay      *nostr.Relay
}

// TODO: use models.Content
type Content struct {
	ID      string
	Content string
	Source  string
}
