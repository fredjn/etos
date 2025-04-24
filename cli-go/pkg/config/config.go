package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	viper *viper.Viper
}

// NewConfig creates a new configuration instance
func NewConfig() *Config {
	v := viper.New()
	config := &Config{
		viper: v,
	}

	// Set default configuration
	config.setDefaults()

	// Set configuration file paths
	config.setConfigPaths()

	return config
}

// setDefaults sets default configuration values
func (c *Config) setDefaults() {
	c.viper.SetDefault("log.level", "info")
	c.viper.SetDefault("plugins.dir", "plugins")
}

// setConfigPaths sets the configuration file search paths
func (c *Config) setConfigPaths() {
	// Set the configuration file name
	c.viper.SetConfigName("etosctl")
	c.viper.SetConfigType("yaml")

	// Add configuration paths
	home, err := os.UserHomeDir()
	if err == nil {
		c.viper.AddConfigPath(filepath.Join(home, ".config", "etos"))
	}
	c.viper.AddConfigPath(".")
}

// Load loads the configuration from files and environment
func (c *Config) Load() error {
	// Load configuration from file
	if err := c.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// Load configuration from environment
	c.viper.SetEnvPrefix("ETOS")
	c.viper.AutomaticEnv()

	return nil
}

// GetString returns a string configuration value
func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}

// GetBool returns a boolean configuration value
func (c *Config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// GetInt returns an integer configuration value
func (c *Config) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// Set sets a configuration value
func (c *Config) Set(key string, value interface{}) {
	c.viper.Set(key, value)
} 