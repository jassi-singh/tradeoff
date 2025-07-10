package config

import (
	"bytes"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"database"`
	Polygon struct {
		APIKey string `mapstructure:"api_key"`
	} `mapstructure:"polygon"`
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	JWT struct {
		Secret     string `mapstructure:"secret"`
		Expiration int64  `mapstructure:"expiration"`
	} `mapstructure:"jwt"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	expandedContent := os.ExpandEnv(string(content))

	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewBufferString(expandedContent)); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
