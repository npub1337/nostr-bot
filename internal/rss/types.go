package rss

type Item struct {
	ID      string
	Title   string
	Link    string
	Content string
	Source  string
}

type Fetcher interface {
	Fetch(urls []string) ([]Item, error)
}
