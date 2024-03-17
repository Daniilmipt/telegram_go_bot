package Config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	ApiKey string `mapstructure:"API_KEY"`
}

// LoadConfig reads config file from specified path and returns Config struct
func LoadConfig(path string) (config Config, err error) {
	// Add path to lookup for config file
	viper.AddConfigPath(path)
	// Set config file name and type
	viper.SetConfigName("conf")
	viper.SetConfigType("env")

	// Try to read config file
	err = viper.ReadInConfig()
	if err != nil {
		_ = fmt.Errorf("do not parse config file: %v", err)
	}

	// Unmarshal config file data to Config struct
	err = viper.Unmarshal(&config)
	if err != nil {
		_ = fmt.Errorf("do not parse config file: %v", err)
	}

	// Return the config struct and the potential error
	return config, err
}
