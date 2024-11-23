package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

type Tweet struct {
	ID   string
	Text string
}

func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Ler e executar o schema
	schemaPath := filepath.Join("sql", "schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) IsTweetPublished(tweetID string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM published_tweets WHERE tweet_id = ?)", tweetID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking tweet existence: %v", err)
		return false
	}
	return exists
}

func (db *DB) MarkTweetAsPublished(tweet Tweet) error {
	_, err := db.Exec(
		"INSERT INTO published_tweets (tweet_id, text) VALUES (?, ?)",
		tweet.ID,
		tweet.Text,
	)
	return err
}
