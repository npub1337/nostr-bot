package rss

type Item struct {
	ID      string
	BotName string
	Title   string
	Link    string
	Content string
	Source  string
}

type Fetcher interface {
	Fetch(name string, urls []string) ([]Item, error)
}
