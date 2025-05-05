package config

import (
	"fmt"
	"log"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	DB     *DBConfig
	Server *ServerConfig
}

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewDBConfig() *DBConfig {
	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		port = 5432 // Use default if parsing fails
	}

	config := &DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "anomaly_detection"),
	}

	log.Printf("Database config: host=%s port=%d user=%s dbname=%s",
		config.Host, config.Port, config.User, config.DBName)

	return config
}

func (c *DBConfig) GetDSN() string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
	log.Printf("Using DSN: %s", dsn)
	return dsn
}
