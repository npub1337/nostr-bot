package nostr

import (
	"github.com/nbd-wtf/go-nostr"
)

type Client struct {
	PrivateKey string
	PublicKey  string
	Relay      *nostr.Relay
}

type Content struct {
	ID      string
	Content string
	Source  string
}
