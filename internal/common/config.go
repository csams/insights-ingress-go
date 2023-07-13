package common

type Config struct {
	*Options
}

type completedConfig struct {
	*Config
}

// CompletedConfig can be constructed only from Config.Complete.
type CompletedConfig struct {
	*completedConfig
}

// NewConfig creates a new configuration object from the given Options.
func NewConfig(o *Options) *Config {
	return &Config{
		Options: o,
	}
}

// Complete ensures the Config has all required values filled.
func (c *Config) Complete() (CompletedConfig, error) {
	return CompletedConfig{&completedConfig{
		c,
	}}, nil
}
