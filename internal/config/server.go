package config

import (
	"fmt"
	"strconv"
)

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int
}

// LoadServerConfig loads configuration from environment variables
func LoadServerConfig() (*ServerConfig, error) {
	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
	}

	serverConfig := &ServerConfig{
		Port: serverPort,
	}

	return serverConfig, nil
}
