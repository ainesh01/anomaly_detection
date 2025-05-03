package config

import "os"

// getEnv returns the value of an environment variable or a default value if it's not set
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
