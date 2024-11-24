CREATE TABLE IF NOT EXISTS content (
    content_id TEXT PRIMARY KEY,
    content TEXT NOT NULL,
    source TEXT NOT NULL, -- "twitter", "rss", etc.
    published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);