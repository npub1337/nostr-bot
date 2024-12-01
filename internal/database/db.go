package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

// TODO: enum for 'status' field

type Content struct {
	ID          string
	Content     string
	Source      string
	Name        string
	LastAttempt time.Time
}

func InitDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	schemaPath := filepath.Join("internal", "database", "schema.sql")
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

func (db *DB) IsContentStored(contentID, botName string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM content WHERE content_id = ? AND bot_name = ?)`
	err := db.QueryRow(query, contentID, botName).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if content is stored: %v", err)
		return false
	}
	return exists
}

func (db *DB) IsContentAlreadyPublished(contentID, botName string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM content WHERE content_id = ? AND bot_name = ? AND published = ?)`
	err := db.QueryRow(query, contentID, botName, 1).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if content is published: %v", err)
		return false
	}
	return exists
}

func (db *DB) InsertRetrievedContent(contentID, content, source, botName string) error {
	query := `INSERT INTO content(content_id, content, source, bot_name) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, contentID, content, source, botName)
	return err
}

func (db *DB) MarkAsPublished(contentID, botName string) error {
	query := `UPDATE content SET published = TRUE WHERE content_id = ? AND bot_name = ?;`
	_, err := db.Exec(query, contentID, botName)
	return err
}

func (db *DB) GetPendingContent() ([]Content, error) {
	query := `
		SELECT content_id, content, source, bot_name
		FROM content 
		WHERE status = 'pending' 
		AND (last_attempt IS NULL OR datetime('now') > datetime(last_attempt, '+5 minutes'))
		AND retry_count < 3
		ORDER BY created_at ASC
		LIMIT 10`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []Content
	for rows.Next() {
		var c Content
		if err := rows.Scan(&c.ID, &c.Content, &c.Source, &c.Name); err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, nil
}

func (db *DB) GetLastPublishedMessage() (Content, error) {
	query := `
		SELECT content_id, content, source, last_attempt
		FROM content 
		WHERE status = 'published' 
		ORDER BY last_attempt DESC
		LIMIT 1`

	row := db.QueryRow(query)

	var c Content
	if err := row.Scan(&c.ID, &c.Content, &c.Source, &c.LastAttempt); err != nil {
		if err == sql.ErrNoRows {
			return c, nil
		}

		return c, err
	}

	return c, nil
}

func (db *DB) UpdateContentStatus(contentID, status string) error {
	query := `
		UPDATE content 
		SET status = ?, 
			last_attempt = datetime('now'),
			retry_count = CASE 
				WHEN ? = 'failed' THEN retry_count + 1 
				ELSE retry_count 
			END
		WHERE content_id = ?`

	_, err := db.Exec(query, status, status, contentID)
	return err
}

// func CloseDB() {
// 	if db != nil {
// 		err := db.Close()
// 		if err != nil {
// 			log.Printf("Error closing the database: %v", err)
// 		}
// 	}
// }
