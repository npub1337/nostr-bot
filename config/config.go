package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RelayURL     string      `yaml:"relay_url"`
	DatabasePath string      `yaml:"database_path"`
	Bots         []BotConfig `yaml:"bots"`
}

type BotConfig struct {
	Name            string   `yaml:"name"`
	NostrPrivateKey string   `yaml:"private_key"`
	RelayURL        string   `yaml:"relay_url,omitempty"`
	RSSFeeds        []string `yaml:"rss_feeds"`
}

type yamlConfig struct {
	Bots []BotConfig `yaml:"bots"`
}

func Load() (*Config, error) {
	config := &Config{
		DatabasePath: "data/content.db",
		Bots:         []BotConfig{},
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

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("error parsing bots config: %w", err)
	}

	return nil
}

func (c *Config) GetRelayURL(botConfig BotConfig) string {
	if botConfig.RelayURL != "" {
		return botConfig.RelayURL
	}
	if c.RelayURL != "" {
		return c.RelayURL
	}
	return "wss://relay.damus.io" //default relay
}
