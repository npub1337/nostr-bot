package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	NostrRelayURL string
	DatabasePath  string
	Bots          []BotConfig
}

type BotConfig struct {
	Name            string   `yaml:"name"`
	NostrPrivateKey string   `yaml:"private_key"`
	RSSFeeds        []string `yaml:"rss_feeds"`
}

type yamlConfig struct {
	Bots []BotConfig `yaml:"bots"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		NostrRelayURL: os.Getenv("NOSTR_RELAY_URL"),
		DatabasePath:  filepath.Join("data", "content.db"),
	}

	if config.NostrRelayURL == "" {
		config.NostrRelayURL = "wss://relay.damus.io"
	}

	if err := config.loadBotsConfig(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) loadBotsConfig() error {
	data, err := os.ReadFile("config/bots.yaml")
	if err != nil {
		return fmt.Errorf("error reading bots config file: %w", err)
	}

	var yamlCfg yamlConfig
	if err := yaml.Unmarshal(data, &yamlCfg); err != nil {
		return fmt.Errorf("error parsing bots config: %w", err)
	}

	c.Bots = yamlCfg.Bots
	return nil
}
