package config

// Default configuration values
const (
	DefaultPort         = "8080"
	DefaultHost         = "0.0.0.0"
	DefaultReadTimeout  = 30
	DefaultWriteTimeout = 30
	DefaultLogLevel     = "info"
	DefaultLogFormat    = "json"
	DefaultConfigPath   = "conf/config.toml"
)

// Environment variable names
const (
	EnvConfigPath = "CONFIG_PATH"
)

// Feature flags
const (
	FeatureSearchEnabled    = true
	FeatureGroupChatEnabled = true
)

// Limits
const (
	DefaultMaxMessageLength = 10000
	DefaultMaxGroupMembers  = 100
)
