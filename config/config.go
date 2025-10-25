package config

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
type Config struct {
	Server ServerConfig
	Log    LogConfig
	Gemini GeminiConfig
	Auth   AuthConfig
}

// ServerConfig stores server specific configuration.
type ServerConfig struct {
	Port int
}

// LogConfig stores logging configuration.
type LogConfig struct {
	Enabled bool
	Path    string
}

// GeminiConfig stores Gemini API specific configuration.
type GeminiConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

// AuthConfig stores authentication configuration.
type AuthConfig struct {
	OpenAIAPIKey string `mapstructure:"openai_api_key"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}