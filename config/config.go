package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config holds the complete application configuration
type Config struct {
	Server   ServerConfig   `toml:"server"`
	Auth     AuthConfig     `toml:"auth"`
	Database DatabaseConfig `toml:"database"`
	Logging  LoggingConfig  `toml:"logging"`
	Features FeaturesConfig `toml:"features"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string `toml:"port"`
	Host         string `toml:"host"`
	ReadTimeout  int    `toml:"read_timeout"`
	WriteTimeout int    `toml:"write_timeout"`
}

// AuthConfig holds authentication configuration with explicit client fields
type AuthConfig struct {
	Client1 ClientAuth `toml:"client1"`
	Client2 ClientAuth `toml:"client2"`
	Client3 ClientAuth `toml:"client3"`
}

// ClientAuth holds client authentication credentials
type ClientAuth struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Mode           string `toml:"mode"`
	MaxConnections int    `toml:"max_connections"`
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level  string `toml:"level"`
	Format string `toml:"format"`
}

// FeaturesConfig holds feature flags and limits
type FeaturesConfig struct {
	EnableSearch     bool `toml:"enable_search"`
	EnableGroupChat  bool `toml:"enable_group_chat"`
	MaxMessageLength int  `toml:"max_message_length"`
	MaxGroupMembers  int  `toml:"max_group_members"`
}

// LoadConfig loads configuration from a TOML file
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = DefaultConfigPath
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Validate required fields
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// NewConfig returns default configuration
func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         DefaultPort,
			Host:         DefaultHost,
			ReadTimeout:  DefaultReadTimeout,
			WriteTimeout: DefaultWriteTimeout,
		},
		Auth: AuthConfig{
			Client1: ClientAuth{
				Username: "user1",
				Password: "password1",
			},
			Client2: ClientAuth{
				Username: "user2",
				Password: "password2",
			},
			Client3: ClientAuth{
				Username: "user3",
				Password: "password3",
			},
		},
		Database: DatabaseConfig{
			Mode:           "memory",
			MaxConnections: 100,
		},
		Logging: LoggingConfig{
			Level:  DefaultLogLevel,
			Format: DefaultLogFormat,
		},
		Features: FeaturesConfig{
			EnableSearch:     FeatureSearchEnabled,
			EnableGroupChat:  FeatureGroupChatEnabled,
			MaxMessageLength: DefaultMaxMessageLength,
			MaxGroupMembers:  DefaultMaxGroupMembers,
		},
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server.port is required")
	}
	// Validate at least one client has credentials
	hasClient := false
	if c.Auth.Client1.Username != "" && c.Auth.Client1.Password != "" {
		hasClient = true
	}
	if c.Auth.Client2.Username != "" && c.Auth.Client2.Password != "" {
		hasClient = true
	}
	if c.Auth.Client3.Username != "" && c.Auth.Client3.Password != "" {
		hasClient = true
	}
	if !hasClient {
		return fmt.Errorf("at least one auth client is required")
	}
	return nil
}

// GetAuthClients returns all authentication clients as a slice for iteration
func (c *Config) GetAuthClients() []ClientAuth {
	clients := []ClientAuth{}
	if c.Auth.Client1.Username != "" && c.Auth.Client1.Password != "" {
		clients = append(clients, c.Auth.Client1)
	}
	if c.Auth.Client2.Username != "" && c.Auth.Client2.Password != "" {
		clients = append(clients, c.Auth.Client2)
	}
	if c.Auth.Client3.Username != "" && c.Auth.Client3.Password != "" {
		clients = append(clients, c.Auth.Client3)
	}
	return clients
}

// ValidateCredentials checks if the provided username and password match any client
func (c *Config) ValidateCredentials(username, password string) bool {
	clients := c.GetAuthClients()
	for _, client := range clients {
		if client.Username == username && client.Password == password {
			return true
		}
	}
	return false
}

// GetPort returns the server port
func (c *Config) GetPort() string {
	return c.Server.Port
}
