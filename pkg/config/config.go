package config

import (
	"advertisement-service/pkg/database"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database database.DBConfig
	Cache    CacheConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type CacheConfig struct {
	Addr     string
	Password string
	DB       string
}

func LoadConfig(configPaths ...string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetDefault("server.port", 8080)

	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	// Unmarshal the config into a Config struct
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return &config, nil
}
