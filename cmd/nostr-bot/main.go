package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"nostr-bot/config"
	"nostr-bot/internal/bot"
	"nostr-bot/internal/database"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Starting Nostr Bot...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	log.Printf("Loaded configuration: %d bots configured", len(cfg.Bots))

	os.MkdirAll(filepath.Dir(cfg.DatabasePath), 0755)

	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	bots := make([]*bot.Bot, 0, len(cfg.Bots))
	for _, botConfig := range cfg.Bots {
		log.Printf("Initializing bot: %s", botConfig.Name)
		b, err := bot.NewBot(
			botConfig.Name,
			botConfig.NostrPrivateKey,
			cfg.RelayURL,
			botConfig.RSSFeeds,
			db,
		)
		if err != nil {
			log.Printf("Error initializing bot %s: %v", botConfig.Name, err)
			continue
		}
		bots = append(bots, b)
		log.Printf("Starting bot: %s", botConfig.Name)
		b.Start()
	}

	log.Printf("All bots started. Running...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	for _, b := range bots {
		b.Stop()
	}
}
