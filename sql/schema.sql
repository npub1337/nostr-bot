CREATE TABLE IF NOT EXISTS published_tweets (
    tweet_id TEXT PRIMARY KEY,
    text TEXT,
    published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
); 